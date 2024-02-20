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
	"github.com/farseer-go/fs/sonyflake"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/file"
	"github.com/farseer-go/utils/http"
	"os"
	"path"
	"strings"
)

// BuildEO 聚合
type BuildEO struct {
	Id                int64               // 主键
	ClusterId         int64               // 集群信息
	BuildNumber       int                 // 构建号
	Status            eumBuildStatus.Enum // 状态
	IsSuccess         bool                // 是否成功
	CreateAt          dateTime.DateTime   // 开始时间
	FinishAt          dateTime.DateTime   // 完成时间
	BuildServerId     int64               // 构建的服务端id（防止生产、开发环境混淆）
	Log               []string            // 构建日志
	Env               EnvVO               // 环境变量
	AppName           string              // 应用名称
	Dockerfile        string              // Dockerfile内容
	WorkflowsAction   ActionVO            // 工作流定义的内容（通过读取WorkflowsYmlPath）
	ShellScript       string              // Shell脚本
	dockerDevice      IDockerDevice
	dockerSwarmDevice IDockerSwarmDevice
	directoryDevice   IDirectoryDevice
	gitDevice         IGitDevice
	kubectlDevice     IKubectlDevice
	copyToDistDevice  ICopyToDistDevice
	logQueue          *LogQueue
	ctx               context.Context
	cancel            context.CancelFunc
	apps              DomainObject
	appGit            GitEO // 应用的源代码
}

func (receiver *BuildEO) IsNil() bool {
	return receiver.Id == 0
}

func (receiver *BuildEO) StartBuild() {
	receiver.ctx, receiver.cancel = context.WithCancel(context.Background())
	receiver.dockerDevice = container.Resolve[IDockerDevice]()
	receiver.dockerSwarmDevice = container.Resolve[IDockerSwarmDevice]()
	receiver.directoryDevice = container.Resolve[IDirectoryDevice]()
	receiver.gitDevice = container.Resolve[IGitDevice]()
	receiver.kubectlDevice = container.Resolve[IKubectlDevice]()
	receiver.copyToDistDevice = container.Resolve[ICopyToDistDevice]()
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

	// 生成Workflows文件
	receiver.checkResult(receiver.GenerateWorkflowsContent())

	// 启动构建系统
	//dockerName := "FOPS-Build-" + strings.NewReplacer(":", "-", ".", "-", "/", "-").Replace(receiver.Env.DockerImage)
	dockerName := "FOPS-Build"
	if !receiver.dockerDevice.ExistsDocker(dockerName) {
		receiver.logQueue.progress <- "启动构建系统：" + receiver.WorkflowsAction.RunsOn
		//args := []string{"-itd", "-v /etc/localtime:/etc/localtime", "-v /var/run/docker.sock:/var/run/docker.sock", "-v /usr/bin/docker:/usr/bin/docker", "-e distRoot=" + DistRoot, "-e gitRoot=" + GitRoot, "-e fopsRoot=" + FopsRoot, "-e npmModulesRoot=" + NpmModulesRoot, "-e kubeRoot=" + KubeRoot, "-e withjson=" + WithJsonPath, "-e dockerfilePath=" + DockerfilePath, "-e dockerIgnorePath=" + DockerIgnorePath, "-e shellRoot=" + ShellRoot, "-e actionsRoot=" + ActionsRoot}
		args := []string{"-itd", "-v /etc/localtime:/etc/localtime", "-v /var/run/docker.sock:/var/run/docker.sock", "-e distRoot=" + DistRoot, "-e gitRoot=" + GitRoot, "-e fopsRoot=" + FopsRoot, "-e npmModulesRoot=" + NpmModulesRoot, "-e kubeRoot=" + KubeRoot, "-e withjson=" + WithJsonPath, "-e dockerfilePath=" + DockerfilePath, "-e dockerIgnorePath=" + DockerIgnorePath, "-e shellRoot=" + ShellRoot, "-e actionsRoot=" + ActionsRoot}
		receiver.checkResult(receiver.dockerDevice.Run(dockerName, "host", receiver.WorkflowsAction.RunsOn, args, true, receiver.Env, receiver.logQueue.progress, receiver.ctx)) // , "-v /var/lib/fops:/var/lib/fops"
	}
	//defer receiver.dockerDevice.Kill(dockerName)
	receiver.logQueue.progress <- "---------------------------------------------------------"

	// 运行step
	for _, step := range receiver.WorkflowsAction.Steps {
		receiver.logQueue.progress <- fmt.Sprintf("执行 %d %s: %s", step.Index, step.Name, step.ActionName)

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
			}

			// 将step文件复制到容器
			receiver.dockerDevice.Copy(dockerName, step.GetActionPath(), step.GetActionPath(), receiver.Env, make(chan string, 100), receiver.ctx)

			// 设置with参数
			step.With["appName"] = receiver.apps.AppName
			step.With["buildId"] = receiver.Env.BuildId
			step.With["buildNumber"] = receiver.Env.BuildNumber
			// 应用的git根目录
			step.With["appAbsolutePath"] = receiver.appGit.GetAbsolutePath()

			// docker
			step.With["dockerImage"] = receiver.Env.DockerImage
			step.With["dockerfilePath"] = receiver.apps.DockerfilePath
			step.With["dockerHub"] = clusterDO.DockerHub
			step.With["dockerUserName"] = clusterDO.DockerUserName
			step.With["dockerUserPwd"] = clusterDO.DockerUserPwd
			step.With["dockerNodeRole"] = receiver.apps.DockerNodeRole
			step.With["dockerNetwork"] = clusterDO.DockerNetwork
			step.With["dockerReplicas"] = receiver.apps.DockerReplicas
			step.With["dockerAdditionalScripts"] = receiver.apps.AdditionalScripts

			// 支持checkout默认拉取应用
			if parse.ToString(step.With["gitHub"]) == "" {
				step.With["gitHub"] = receiver.appGit.Hub
				step.With["gitBranch"] = receiver.appGit.Branch
				step.With["gitUserName"] = receiver.appGit.UserName
				step.With["gitUserPwd"] = receiver.appGit.UserPwd
				step.With["gitPath"] = receiver.appGit.Dir
			}

			// 生成with.json文件，并复制到容器
			file.Delete(WithJsonPath)
			withContent, _ := json.Marshal(step.With)
			file.WriteByte(WithJsonPath, withContent)
			receiver.dockerDevice.Copy(dockerName, WithJsonPath, WithJsonPath, receiver.Env, make(chan string, 100), receiver.ctx)

			// 执行 docker exec FOPS-Build-hub-fsgit-cc-fops-130 echo aaa
			receiver.checkResult(receiver.dockerDevice.Execute(dockerName, step.GetActionPath(), receiver.Env, receiver.logQueue.progress, receiver.ctx))
		}

		// 运行脚本
		if len(step.Run) > 0 {
			shellScript := collections.NewList[string]()
			shellScript.Add("export PATH=$PATH:/usr/local/go/bin")
			shellScript.Add("go env -w GO111MODULE=on && go env -w GOPROXY=https://goproxy.cn,direct")
			shellScript.Add("cd " + DistRoot + receiver.appGit.GetRelativePath())
			shellScript.AddArray(step.Run)
			shellScript.Add("")
			shellPath := fmt.Sprintf("%s%d-%d.sh", ShellRoot, receiver.Env.BuildNumber, step.Index)
			file.WriteString(shellPath, shellScript.ToString("\n"))
			receiver.dockerDevice.Copy(dockerName, shellPath, shellPath, receiver.Env, make(chan string, 100), receiver.ctx)

			receiver.checkResult(exec.RunShell("docker exec "+dockerName+" /bin/sh -x "+shellPath, receiver.logQueue.progress, receiver.Env.ToMap(), DistRoot, false) == 0)
			//receiver.checkResult(receiver.dockerDevice.Execute(dockerName, step.Run, receiver.Env, receiver.logQueue.progress, receiver.ctx))
		}
		receiver.logQueue.progress <- "---------------------------------------------------------"
	}

	//// 打印环境变量
	//receiver.Env.Print(receiver.logQueue.progress)
	//
	//// 前置检查
	//receiver.directoryDevice.Check(receiver.logQueue.progress)
	//
	//// 拉取主仓库及依赖仓库
	//receiver.checkResult(receiver.gitDevice.CloneOrPullAndDependent(receiver.getGits(), receiver.logQueue.progress, receiver.ctx))
	//
	//// 登陆镜像仓库(先登陆，如果失败了，后则面也不需要编译、打包了)
	//receiver.checkResult(receiver.dockerDevice.Login(clusterDO.DockerHub, clusterDO.DockerUserName, clusterDO.DockerUserPwd, receiver.logQueue.progress, receiver.Env, receiver.ctx))
	//
	//// 将需要打包的源代码，复制到dist目录
	//receiver.copyToDistDevice.Copy(receiver.getGits(), receiver.Env, receiver.logQueue.progress)
	//
	//// 生成Dockerfile文件
	//receiver.checkResult(receiver.GenerateDockerfileContent())
	//receiver.dockerDevice.CreateDockerfile(receiver.AppName, receiver.Dockerfile, receiver.ctx)
	//
	//// docker打包
	//receiver.checkResult(receiver.dockerDevice.Build(receiver.Env, receiver.logQueue.progress, receiver.ctx))
	//
	//// docker上传
	//receiver.checkResult(receiver.dockerDevice.Push(receiver.Env, receiver.logQueue.progress, receiver.ctx))
	//
	//// 首次创建还是更新镜像
	//if receiver.dockerSwarmDevice.ExistsDocker(clusterDO, receiver.AppName) {
	//	// 更新镜像
	//	receiver.checkResult(receiver.dockerSwarmDevice.SetImages(clusterDO, receiver.AppName, receiver.Env.DockerImage, receiver.logQueue.progress, receiver.ctx))
	//} else {
	//	// 创建容器服务
	//	receiver.checkResult(receiver.dockerSwarmDevice.CreateService(receiver.AppName, receiver.apps.DockerNodeRole, receiver.apps.AdditionalScripts, clusterDO.DockerNetwork, receiver.apps.DockerReplicas, receiver.Env.DockerImage, receiver.logQueue.progress, receiver.ctx))
	//}
	receiver.success()
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

	// 发布事件
	event.BuildFinishedEvent{AppName: receiver.AppName, BuildId: receiver.Id, ClusterId: receiver.ClusterId, IsSuccess: true}.PublishEvent()

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

// GenerateDockerfileContent 生成Dockerfile
func (receiver *BuildEO) GenerateDockerfileContent() bool {
	// 为空时，读取应用git文件中的Dockerfile文件
	if receiver.Dockerfile == "" {
		// 如果没有自定义，则使用应用仓库根目录的Dockerfile文件
		if receiver.apps.DockerfilePath == "" {
			receiver.apps.DockerfilePath = "Dockerfile"
		} else {
			// 自定义Dockerfile路径
			if strings.HasPrefix(receiver.apps.DockerfilePath, "/") {
				receiver.apps.DockerfilePath = receiver.apps.DockerfilePath[1:]
			} else if strings.HasPrefix(receiver.apps.DockerfilePath, "./") {
				receiver.apps.DockerfilePath = receiver.apps.DockerfilePath[2:]
			}
		}
		receiver.apps.DockerfilePath = receiver.appGit.GetAbsolutePath() + receiver.apps.DockerfilePath
		receiver.Dockerfile = file.ReadString(receiver.apps.DockerfilePath)
		if receiver.Dockerfile == "" {
			receiver.logQueue.progress <- "Dockerfile没有定义。"
			return false
		}
		return true
	}

	// 替换项目名称
	receiver.Dockerfile = strings.ReplaceAll(receiver.Dockerfile, "${app_name}", receiver.AppName)
	receiver.Dockerfile = strings.ReplaceAll(receiver.Dockerfile, "${git_name}", receiver.Env.GitName)
	//receiver.Dockerfile = strings.ReplaceAll(receiver.Dockerfile, "${domain}", do.Project.Domain)
	//receiver.Dockerfile = strings.ReplaceAll(receiver.Dockerfile, "${entry_point}", do.Project.EntryPoint)
	//receiver.Dockerfile = strings.ReplaceAll(receiver.Dockerfile, "${entry_port}", strconv.Itoa(do.Project.EntryPort))
	//receiver.Dockerfile = strings.ReplaceAll(receiver.Dockerfile, "${project_path}", strings.TrimPrefix(do.Project.Path, "/"))
	return true
}

// GenerateWorkflowsContent 生成Workflows
func (receiver *BuildEO) GenerateWorkflowsContent() bool {
	// 如果没有自定义，则使用应用仓库根目录的.fops/workflows/build.yml文件
	if receiver.apps.WorkflowsYmlPath == "" {
		receiver.apps.WorkflowsYmlPath = receiver.appGit.GetRawContent(".fops/workflows/build.yml")
		// 自定义WorkflowsYmlPath路径
	} else if !strings.HasPrefix(receiver.apps.WorkflowsYmlPath, "http://") && !strings.HasPrefix(receiver.apps.WorkflowsYmlPath, "https://") {
		receiver.apps.WorkflowsYmlPath = receiver.appGit.GetRawContent(receiver.apps.WorkflowsYmlPath)
	}

	receiver.apps.WorkflowsYmlPath = receiver.apps.WorkflowsYmlPath + "?" + parse.ToString(sonyflake.GenerateId())
	receiver.logQueue.progress <- "加载工作流文件：" + receiver.apps.WorkflowsYmlPath

	// 通过http读取工作流定义的内容
	var err error
	receiver.WorkflowsAction, err = LoadWorkflows(receiver.apps.WorkflowsYmlPath, receiver.AppName, receiver.Env.GitName)
	if err != nil {
		receiver.logQueue.progress <- err.Error()
		return false
	}

	receiver.WorkflowsAction.Steps = append([]stepVO{
		{
			Index:             1,
			Name:              "初始化环境",
			ActionName:        "clear",
			ActionVer:         "v1",
			ActionDownloadUrl: "https://github.com/farseers/FOPS-Actions/releases/download/v1/clear",
			RepositoryName:    "farseers/FOPS-Actions",
			With:              make(map[string]any),
		},
	}, receiver.WorkflowsAction.Steps...)

	receiver.logQueue.progress <- "读取到工作流文件：" + receiver.WorkflowsAction.Name
	return true
}
