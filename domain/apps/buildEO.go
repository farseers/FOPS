package apps

import (
	"context"
	"encoding/json"
	"fmt"
	"fops/domain/_/eumBuildStatus"
	"fops/domain/apps/event"
	"fops/domain/cluster"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/utils/file"
	"github.com/farseer-go/utils/http"
	"os"
	"path"
	"strings"
)

// BuildEO 聚合
type BuildEO struct {
	Id              int64               // 主键
	ClusterId       int64               // 集群信息
	BuildNumber     int                 // 构建号
	Status          eumBuildStatus.Enum // 状态
	IsSuccess       bool                // 是否成功
	CreateAt        dateTime.DateTime   // 开始时间
	FinishAt        dateTime.DateTime   // 完成时间
	BuildServerId   int64               // 构建的服务端id（防止生产、开发环境混淆）
	Env             EnvVO               // 环境变量
	AppName         string              // 应用名称
	WorkflowsName   string              // 工作流名称（文件的名称）
	WorkflowsAction ActionVO            // 工作流定义的内容（通过读取WorkflowsYmlPath）
	dockerDevice    IDockerDevice
	gitDevice       IGitDevice
	logQueue        *LogQueue
	ctx             context.Context
	cancel          context.CancelFunc
	apps            DomainObject
	appGit          GitEO // 应用的源代码
}

func (receiver *BuildEO) IsNil() bool {
	return receiver.Id == 0
}

func (receiver *BuildEO) StartBuild() {
	receiver.ctx, receiver.cancel = context.WithCancel(context.Background())
	receiver.dockerDevice = container.Resolve[IDockerDevice]()
	receiver.gitDevice = container.Resolve[IGitDevice]()
	receiver.logQueue = NewLogQueue(receiver.Id)

	appsRepository := container.Resolve[Repository]()
	// 应用
	receiver.apps = appsRepository.ToEntity(receiver.AppName)
	receiver.appGit = appsRepository.ToGitEntity(receiver.apps.AppGit)

	// 集群
	clusterRepository := container.Resolve[cluster.Repository]()
	clusterDO := clusterRepository.ToEntity(receiver.ClusterId)

	// 定义环境变量
	var projectGitRoot = receiver.appGit.GetAbsolutePath()
	var dockerHub = receiver.dockerDevice.GetDockerHub(clusterDO.DockerHub)
	var dockerImage = receiver.dockerDevice.GetDockerImage(clusterDO.DockerHub, receiver.AppName, receiver.BuildNumber)
	var gitName = receiver.appGit.GetName()
	receiver.GenerateEnv(projectGitRoot, dockerHub, dockerImage, gitName)

	defer receiver.catch()

	// 设置with参数
	sysWith := map[string]any{
		"appName":     receiver.apps.AppName,
		"buildId":     receiver.Env.BuildId,
		"buildNumber": receiver.Env.BuildNumber,
		// 应用的git根目录
		"appAbsolutePath": receiver.appGit.GetAbsolutePath(),
		// docker
		"dockerImage":             receiver.Env.DockerImage,
		"dockerfilePath":          receiver.apps.DockerfilePath,
		"dockerHub":               clusterDO.DockerHub,
		"dockerUserName":          clusterDO.DockerUserName,
		"dockerUserPwd":           clusterDO.DockerUserPwd,
		"dockerNodeRole":          receiver.apps.DockerNodeRole,
		"dockerNetwork":           clusterDO.DockerNetwork,
		"dockerReplicas":          receiver.apps.DockerReplicas,
		"dockerAdditionalScripts": receiver.apps.AdditionalScripts,
		"clusterId":               receiver.ClusterId,
	}

	// 把fops、fschedule版本写入到系统参数sysWith
	//sysWith["fops.ver"] = container.Resolve[Repository]().ToEntity("fops").DockerVer
	//sysWith["fschedule.ver"] = container.Resolve[Repository]().ToEntity("fschedule").DockerVer

	// 生成Workflows文件
	receiver.checkResult(receiver.GenerateWorkflowsContent(sysWith))

	// 启动构建系统
	dockerName := "FOPS-Build"
	if !receiver.dockerDevice.ExistsDocker(dockerName) {
		receiver.logQueue.progress <- "启动构建系统：" + receiver.WorkflowsAction.RunsOn
		args := []string{"-itd", "-v /etc/localtime:/etc/localtime", "-v /var/run/docker.sock:/var/run/docker.sock", "-e distRoot=" + DistRoot, "-e gitRoot=" + GitRoot, "-e fopsRoot=" + FopsRoot, "-e npmModulesRoot=" + NpmModulesRoot, "-e kubeRoot=" + KubeRoot, "-e withjson=" + WithJsonPath, "-e dockerfilePath=" + DockerfilePath, "-e dockerIgnorePath=" + DockerIgnorePath, "-e shellRoot=" + ShellRoot, "-e actionsRoot=" + ActionsRoot}
		receiver.checkResult(receiver.dockerDevice.Run(dockerName, "host", receiver.WorkflowsAction.RunsOn, args, true, receiver.Env, receiver.logQueue.progress, receiver.ctx)) // , "-v /var/lib/fops:/var/lib/fops"
	}
	//defer receiver.dockerDevice.Kill(dockerName)

	// 设置镜像的代理
	if receiver.WorkflowsAction.Proxy != "" {
		receiver.WorkflowsAction.Env["HTTP_PROXY"] = "http://" + receiver.WorkflowsAction.Proxy
		receiver.WorkflowsAction.Env["HTTPS_PROXY"] = "http://" + receiver.WorkflowsAction.Proxy
	}

	// 加载环境变量提示
	if len(receiver.WorkflowsAction.Env) > 0 {
		receiver.logQueue.progress <- "加载环境变量："
		for k, v := range receiver.WorkflowsAction.Env {
			receiver.logQueue.progress <- fmt.Sprintf("%s=%s", k, v)
		}
	}
	receiver.logQueue.progress <- "---------------------------------------------------------"

	// 运行step
	for index, step := range receiver.WorkflowsAction.Steps {
		receiver.logQueue.progress <- fmt.Sprintf("执行 %d %s: %s", index+1, step.ActionName, step.Name)

		// 使用action程序，需要判断是否要下载
		if step.ActionName != "" {
			if !file.IsExists(step.GetActionPath()) {
				receiver.logQueue.progress <- fmt.Sprintf("下载 %s", step.ActionDownloadUrl)
				// 先创建目录
				file.CreateDir766(path.Dir(step.GetActionPath()))
				// 下载文件
				if err := http.Download(step.ActionDownloadUrl, step.GetActionPath(), 0, configure.GetString("Fops.GitAgent")); err != nil {
					receiver.logQueue.progress <- fmt.Sprintf("下载action %s 时发生错误：%s", step.ActionDownloadUrl, err.Error())
					receiver.checkResult(false)
				}
				_ = os.Chmod(step.GetActionPath(), 777)
				receiver.logQueue.progress <- "下载完成"
			}

			// 将action文件复制到容器
			receiver.dockerDevice.Copy(dockerName, step.GetActionPath(), step.GetActionPath(), receiver.Env, make(chan string, 100), receiver.ctx)

			gits := receiver.getGits()
			// 支持checkout默认拉取应用
			if parse.ToString(step.With["gitHub"]) != "" {
				gits = append(gits, GitEO{
					Hub:      parse.ToString(step.With["gitHub"]),
					Branch:   parse.ToString(step.With["gitBranch"]),
					UserName: parse.ToString(step.With["gitUserName"]),
					UserPwd:  parse.ToString(step.With["gitUserPwd"]),
					Path:     parse.ToString(step.With["gitPath"]),
				})
			}
			step.With["gits"] = gits

			// 生成with.json文件，并复制到容器
			file.Delete(WithJsonPath)
			withContent, _ := json.Marshal(step.With)
			file.WriteByte(WithJsonPath, withContent)
			receiver.dockerDevice.Copy(dockerName, WithJsonPath, WithJsonPath, receiver.Env, make(chan string, 100), receiver.ctx)

			// 执行 docker exec FOPS-Build-hub-fsgit-cc-fops-130 echo aaa
			receiver.checkResult(receiver.dockerDevice.Execute(dockerName, step.GetActionPath(), receiver.WorkflowsAction.Env, receiver.logQueue.progress, receiver.ctx))

			switch step.ActionName {
			case "checkout":
				event.GitCloneOrPulledEvent{GitId: receiver.appGit.Id}.PublishEvent()
			case "dockerPush":
				// 上传成功后，需要更新项目中的镜像版本属性
				event.DockerPushedEvent{BuildNumber: parse.ToInt(step.With["buildNumber"]), AppName: parse.ToString(step.With["appName"]), ImageName: parse.ToString(step.With["dockerImage"])}.PublishEvent()
			}
		}

		// 运行脚本
		if len(step.Run) > 0 {
			shellScript := collections.NewList[string]()
			shellScript.Add("source /etc/profile")
			//shellScript.Add("go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct")
			shellScript.Add("mkdir -p " + DistRoot + receiver.appGit.GetRelativePath())
			shellScript.Add("cd " + DistRoot + receiver.appGit.GetRelativePath())
			shellScript.AddArray(step.Run)
			shellScript.Add("")
			script := shellScript.ToString("\n")
			// 支持参数化脚本
			for k, v := range step.With {
				script = strings.ReplaceAll(script, "{{"+k+"}}", parse.ToString(v))
			}
			shellPath := fmt.Sprintf("%s%d-%d.sh", ShellRoot, receiver.Env.BuildNumber, index+1)
			file.WriteString(shellPath, script)
			receiver.dockerDevice.Copy(dockerName, shellPath, shellPath, receiver.Env, make(chan string, 100), receiver.ctx)

			receiver.checkResult(receiver.dockerDevice.Execute(dockerName, "/bin/sh -x "+shellPath, receiver.WorkflowsAction.Env, receiver.logQueue.progress, receiver.ctx))
			//receiver.checkResult(exec.RunShell("docker exec "+dockerName+" /bin/sh -x "+shellPath, receiver.logQueue.progress, receiver.Env.ToMap(), DistRoot, false) == 0)
		}
		receiver.logQueue.progress <- "---------------------------------------------------------"
	}

	receiver.success()
}

// GenerateWorkflowsContent 生成Workflows
func (receiver *BuildEO) GenerateWorkflowsContent(sysWith map[string]any) bool {
	// 更新工作流文件到本地
	receiver.gitDevice.PullWorkflows(receiver.apps.GetWorkflowsRoot(), receiver.appGit.Branch, receiver.appGit.GetAuthHub(), receiver.logQueue.progress)

	// 通过http读取工作流定义的内容
	var err error
	receiver.WorkflowsAction, err = LoadWorkflows(receiver.apps.GetWorkflowsDir()+receiver.WorkflowsName+".yml", receiver.AppName, receiver.Env.GitName, sysWith)
	if err != nil {
		receiver.logQueue.progress <- err.Error()
		return false
	}

	//if gitAgent := configure.GetString("Fops.GitAgent"); gitAgent != "" {
	//	receiver.WorkflowsAction.Steps = append([]stepVO{
	//		{
	//			Name:              "开启Git代理",
	//			ActionName:        "gitProxy",
	//			ActionVer:         "v1",
	//			ActionDownloadUrl: "https://github.com/farseers/FOPS-Actions/releases/download/v1/gitProxy",
	//			RepositoryName:    "farseers/FOPS-Actions",
	//			With:              map[string]any{"proxy": gitAgent},
	//		},
	//	}, receiver.WorkflowsAction.Steps...)
	//}

	receiver.WorkflowsAction.Steps = append([]stepVO{
		{
			Name:              "初始化环境",
			ActionName:        "clear",
			ActionVer:         "v1",
			ActionDownloadUrl: "https://github.com/farseers/FOPS-Actions/releases/download/v1/clear",
			RepositoryName:    "farseers/FOPS-Actions",
			With:              make(map[string]any),
		},
	}, receiver.WorkflowsAction.Steps...)

	receiver.logQueue.progress <- "读取到工作流文件：" + receiver.WorkflowsAction.Name

	// 将全局参数 覆盖到 系统参数
	for k, v := range receiver.WorkflowsAction.With {
		switch v.(type) {
		case string: // 字符串的类型，才需要替换
			for sysKey, sysVal := range sysWith {
				v = strings.ReplaceAll(parse.ToString(v), "{{"+sysKey+"}}", parse.ToString(sysVal))
			}
		}
		sysWith[k] = v
	}

	// 系统参数替换到step.with变量
	for _, step := range receiver.WorkflowsAction.Steps {
		// 替换参数变量
		for k, v := range step.With {
			switch v.(type) {
			case string: // 字符串的类型，才需要替换
				for sysKey, sysVal := range sysWith {
					step.With[k] = strings.ReplaceAll(parse.ToString(step.With[k]), "{{"+sysKey+"}}", parse.ToString(sysVal))
				}
			}
		}

		// 系统参数 合并到 step.with
		for k, v := range sysWith {
			// 系统参数 和 自定义参数 同时有的话，忽略
			if _, exists := step.With[k]; !exists {
				step.With[k] = v
			}
		}
	}

	return true
}

func (receiver *BuildEO) catch() {
	if err := recover(); err != nil {
		receiver.cancel()
		var msg string
		switch e := err.(type) {
		case string:
			msg = e
		case exception.RefuseException:
			msg = e.Message
		}
		if msg != "exit" {
			receiver.logQueue.progress <- msg
		}
		receiver.fail()
	}
	receiver.logQueue.Close()
}

// CheckResult 检查结构
func (receiver *BuildEO) checkResult(result bool) {
	status := container.Resolve[Repository]().GetStatus(receiver.Id)
	if status == eumBuildStatus.Finish {
		exception.ThrowRefuseException("手动取消，退出构建。")
	}

	if !result {
		exception.ThrowRefuseException("exit")
	}
}

// 设置任务失败
func (receiver *BuildEO) fail() {
	receiver.logQueue.progress <- "---------------------------------------------------------"
	receiver.logQueue.progress <- "执行失败，退出构建。"

	// 发布事件
	event.BuildFinishedEvent{AppName: receiver.AppName, BuildId: receiver.Id, ClusterId: receiver.ClusterId, IsSuccess: false}.PublishEvent()

	container.Resolve[Repository]().SetCancel(receiver.Id)
}

// 设置任务成功
func (receiver *BuildEO) success() {
	receiver.logQueue.progress <- "---------------------------------------------------------"
	receiver.logQueue.progress <- "构建完成。"

	// 包含dockerswarmUpdateVer，才要发布通知
	if collections.NewList(receiver.WorkflowsAction.Steps...).Where(func(item stepVO) bool {
		return item.ActionName == "dockerswarmUpdateVer"
	}).Any() {
		receiver.logQueue.progress <- "更新部署版本。"
		// 发布事件
		event.BuildFinishedEvent{AppName: receiver.AppName, BuildId: receiver.Id, ClusterId: receiver.ClusterId, IsSuccess: true}.PublishEvent()
	}

	container.Resolve[Repository]().SetSuccess(receiver.Id)
}

// 得到所有Git
func (receiver *BuildEO) getGits() []GitEO {
	var gits []GitEO
	if !receiver.appGit.IsNil() {
		gits = append(gits, receiver.appGit)
	}
	// 依赖的框架
	frameworkGits := container.Resolve[Repository]().ToGitList(receiver.apps.FrameworkGits)
	if frameworkGits.Count() > 0 {
		gits = append(gits, frameworkGits.ToArray()...)
	}
	return gits
}

// GenerateEnv 生成环境变量
func (receiver *BuildEO) GenerateEnv(projectGitRoot string, dockerHub string, dockerImage string, gitName string) {
	receiver.Env = EnvVO{
		BuildId:     receiver.Id,
		BuildNumber: receiver.BuildNumber,
		AppName:     receiver.AppName,
		DockerHub:   dockerHub,
		DockerImage: dockerImage,
		AppGitRoot:  projectGitRoot,
		GitHub:      receiver.appGit.Hub,
		GitName:     gitName,
	}
}
