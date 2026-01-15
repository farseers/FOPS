package apps

import (
	"context"

	"fmt"
	"fops/domain/_/eumBuildStatus"
	"fops/domain/_/eumBuildType"
	"fops/domain/apps/event"
	"fops/domain/appsBranch"
	"fops/domain/cluster"
	"fops/domain/monitor"
	"os"
	"path"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/queue"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/file"
	"github.com/farseer-go/utils/http"
)

// BuildEO 聚合
type BuildEO struct {
	Id              int64               // 主键
	ClusterId       int64               // 集群信息
	BuildNumber     int                 // 构建号
	Status          eumBuildStatus.Enum // 状态
	BuildType       eumBuildType.Enum   // 构建类型
	IsSuccess       bool                // 是否成功
	CreateAt        dateTime.DateTime   // 开始时间
	FinishAt        dateTime.DateTime   // 完成时间
	BuildServerId   int64               // 构建的服务端id（防止生产、开发环境混淆）
	Env             EnvVO               // 环境变量
	AppName         string              // 应用名称
	WorkflowsName   string              // 工作流名称（文件的名称）
	BranchName      string              // 分支名称
	DockerImage     string              // Docker镜像
	WorkflowsAction ActionVO            // 工作流定义的内容（通过读取WorkflowsYmlPath）
	dockerDevice    IDockerDevice
	gitDevice       IGitDevice
	logQueue        *LogQueue
	ctx             context.Context
	cancel          context.CancelFunc
	apps            DomainObject
	appGit          GitEO // 应用的源代码
	dockerClient    *docker.Client
	fopsBuildName   string // 构建的容器名称
}

func (receiver *BuildEO) IsNil() bool {
	return receiver.Id == 0
}

func (receiver *BuildEO) StartBuild() {
	switch receiver.BuildType {
	case eumBuildType.Manual:
		receiver.fopsBuildName = "FOPS-Build"
	case eumBuildType.Auto:
		receiver.fopsBuildName = "FOPS-AutoBuild"
	}

	receiver.ctx, receiver.cancel = context.WithCancel(context.Background())
	receiver.dockerClient = docker.NewClient()
	receiver.dockerDevice = container.Resolve[IDockerDevice]()
	receiver.gitDevice = container.Resolve[IGitDevice]()
	receiver.logQueue = NewLogQueue(receiver.Id)

	appsRepository := container.Resolve[Repository]()

	// 应用
	receiver.apps = appsRepository.ToEntity(receiver.AppName)
	receiver.appGit = appsRepository.ToGitEntity(receiver.apps.AppGit)

	// 开启异步监控状态
	go receiver.WatchStatus()
	defer receiver.catch()

	// 生成Workflows（并更新集群ID）
	receiver.checkResult(receiver.GenerateWorkflowsContent())
	// 集群
	clusterRepository := container.Resolve[cluster.Repository]()
	clusterDO := clusterRepository.ToEntity(receiver.ClusterId)
	if clusterDO.IsNil() {
		receiver.logQueue.progress <- fmt.Sprintf("集群不存在：%d", receiver.ClusterId)
		receiver.checkResult(false)
	}
	if clusterDO.DockerNetwork == "" {
		receiver.logQueue.progress <- "集群的容器网络未配置"
		receiver.checkResult(false)
	}

	// 定义环境变量
	var projectGitRoot = receiver.appGit.GetAbsolutePath()
	receiver.DockerImage = receiver.dockerDevice.GetDockerImage(clusterDO.DockerHub, receiver.AppName, receiver.BuildNumber)
	var dockerHub = receiver.dockerDevice.GetDockerHub(clusterDO.DockerHub)
	receiver.GenerateEnv(projectGitRoot, dockerHub, receiver.DockerImage, receiver.appGit.GetName())
	// 更新集群ID、镜像（需要先生成好Env)
	appsRepository.UpdateBuilding(receiver.Id, receiver.Env)

	receiver.ReplaceSysWith(map[string]any{
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
		"isLocal":                 clusterDO.IsLocal,
		"fopsAddr":                clusterDO.FopsAddr,
		"fScheduleAddr":           clusterDO.FScheduleAddr,
		"limitCpus":               receiver.apps.LimitCpus,
		"limitMemory":             receiver.apps.LimitMemory,
	})

	// 设置镜像的代理
	if receiver.WorkflowsAction.With["proxy"] != "" {
		receiver.WorkflowsAction.Env["HTTP_PROXY"] = parse.ToString(receiver.WorkflowsAction.With["proxy"])
		receiver.WorkflowsAction.Env["HTTPS_PROXY"] = parse.ToString(receiver.WorkflowsAction.With["proxy"])
	}

	// 加载环境变量提示
	if len(receiver.WorkflowsAction.Env) > 0 {
		receiver.logQueue.progress <- "构建环境变量："
		for k, v := range receiver.WorkflowsAction.Env {
			receiver.logQueue.progress <- fmt.Sprintf("%s=%s", k, v)
		}
	}

	if clusterDO.IsLocal {
		receiver.logQueue.progress <- "部署到本地集群"
	} else {
		receiver.logQueue.progress <- "部署到远程集群：" + clusterDO.FopsAddr
	}

	// 启动构建系统
	if !receiver.dockerClient.Container.Exists(receiver.fopsBuildName) {
		receiver.logQueue.progress <- "启动构建系统：" + receiver.WorkflowsAction.RunsOn
		args := []string{"-itd", "-v /etc/localtime:/etc/localtime", "-v /var/run/docker.sock:/var/run/docker.sock", "-e distRoot=" + DistRoot, "-e gitRoot=" + GitRoot, "-e fopsRoot=" + FopsRoot, "-e npmModulesRoot=" + NpmModulesRoot, "-e kubeRoot=" + KubeRoot, "-e withjson=" + WithJsonPath, "-e dockerfilePath=" + DockerfilePath, "-e dockerIgnorePath=" + DockerIgnorePath, "-e shellRoot=" + ShellRoot, "-e actionsRoot=" + ActionsRoot}
		// 添加自定义的挂载
		builderArgs := configure.GetSlice("Fops.Builder")
		args = append(args, builderArgs...)
		lstResult, wait := receiver.dockerClient.Container.Run(receiver.fopsBuildName, "host", receiver.WorkflowsAction.RunsOn, args, true, receiver.Env.ToMap(), receiver.ctx) // , "-v /var/lib/fops:/var/lib/fops"
		if wait() != 0 {
			for result := range lstResult {
				receiver.logQueue.progress <- result
			}
			receiver.checkResult(false)
		}
	}
	//defer receiver.dockerClient.Container.Kill(receiver.fopsBuildName)

	// 获取外网IP
	if extranetIpUrl := configure.GetString("Fops.ExtranetIpUrl"); extranetIpUrl != "" {
		receiver.ctx, receiver.cancel = context.WithCancel(context.Background())
		lstResult, wait := receiver.dockerClient.Container.Exec(receiver.fopsBuildName, fmt.Sprintf("echo '公网IP：$(curl -s %s)'", extranetIpUrl), nil, receiver.ctx)
		exec.SaveToChan(receiver.logQueue.progress, lstResult, wait)
	}
	receiver.logQueue.progress <- "---------------------------------------------------------"

	gits := receiver.getGits()

	// 运行step
	for index, step := range receiver.WorkflowsAction.Steps {
		// 执行Action
		receiver.runStep(index, step, gits)
		// 针对特殊的Action，附加额外逻辑
		switch step.ActionName {
		case "checkout":
			// 得到应用的CommitId
			receiver.getCommitId()
			receiver.getSha256sum()
			// 直接使用缓存
			if receiver.useCache(index, gits) {
				receiver.success()
				return
			}

		case "dockerPush": // 上传成功后，需要更新项目中的镜像版本属性
			//event.DockerPushedEvent{BuildNumber: parse.ToInt(step.With["buildNumber"]), AppName: parse.ToString(step.With["appName"]), ImageName: parse.ToString(step.With["dockerImage"])}.PublishEvent()
		case "dockerBuild": // 镜像打包成功后，需要更新到Git分支中，用于后续的缓存使用
			container.Resolve[appsBranch.Repository]().UpdateDockerImage(receiver.AppName, receiver.Env.CommitId, receiver.Env.DockerImage, receiver.Env.Sha256sum)
		}
	}

	receiver.success()
}

// 执行Action
func (receiver *BuildEO) runStep(index int, step stepVO, gits collections.List[GitEO]) {
	receiver.logQueue.progress <- fmt.Sprintf("执行 %d %s: %s", index+1, step.ActionName, step.Name)
	// 使用action程序，需要判断是否要下载
	if step.ActionName != "" {
		receiver.ctx, receiver.cancel = context.WithCancel(context.Background())
		if !file.IsExists(step.GetActionPath()) {
			receiver.logQueue.progress <- fmt.Sprintf("下载 %s", step.ActionDownloadUrl)
			// 先创建目录
			file.CreateDir766(path.Dir(step.GetActionPath()))
			// 下载文件
			if _, err := http.Download(step.ActionDownloadUrl, step.GetActionPath(), nil, 0, configure.GetString("Fops.Proxy")); err != nil {
				receiver.logQueue.progress <- fmt.Sprintf("下载action %s 时发生错误：%s", step.ActionDownloadUrl, err.Error())
				receiver.checkResult(false)
			}
			_ = os.Chmod(step.GetActionPath(), 777)
			receiver.logQueue.progress <- "下载完成"
		}

		// 将action文件复制到容器
		_, wait := receiver.dockerClient.Container.Cp(receiver.fopsBuildName, step.GetActionPath(), step.GetActionPath(), receiver.ctx)
		wait()

		// 支持checkout默认拉取应用
		if step.ActionName == "checkout" {
			// 自定义了地址
			if parse.ToString(step.With["gitHub"]) != "" {
				gits.Add(GitEO{
					Hub:      parse.ToString(step.With["gitHub"]),
					Branch:   parse.ToString(step.With["gitBranch"]),
					UserName: parse.ToString(step.With["gitUserName"]),
					UserPwd:  parse.ToString(step.With["gitUserPwd"]),
					Path:     parse.ToString(step.With["gitPath"]),
				})
			}
			// 修改了应用的分支
			if branch := parse.ToString(step.With["branch"]); branch != "" || receiver.BranchName != "" {
				appGit := gits.Find(func(item *GitEO) bool {
					return item.IsApp
				})
				// UI中传入
				if receiver.BranchName != "" {
					appGit.Branch = receiver.BranchName
				} else {
					// 使用工作流定义的分支
					appGit.Branch = branch
					receiver.BranchName = branch
					receiver.Env.BranchName = branch
					container.Resolve[Repository]().UpdateBuilding(receiver.Id, receiver.Env)
				}
			}
		}
		step.With["gits"] = gits

		// 生成with.json文件，并复制到容器
		file.Delete(WithJsonPath)
		withContent, _ := snc.Marshal(step.With)
		file.WriteByte(WithJsonPath, withContent)
		_, wait = receiver.dockerClient.Container.Cp(receiver.fopsBuildName, WithJsonPath, WithJsonPath, receiver.ctx)
		wait()

		// 设置超时
		receiver.ctx, receiver.cancel = context.WithTimeout(receiver.ctx, time.Duration(step.Timeout)*time.Minute)

		// 执行 docker exec FOPS-Build /bin/bash -c "action"
		lstResult, wait := receiver.dockerClient.Container.Exec(receiver.fopsBuildName, step.GetActionPath(), receiver.WorkflowsAction.Env, receiver.ctx)
		code := exec.SaveToChan(receiver.logQueue.progress, lstResult, wait)
		receiver.checkResult(code == 0)
	}

	// 运行脚本
	if len(step.Run) > 0 {
		receiver.ctx, receiver.cancel = context.WithCancel(context.Background())
		shellScript := collections.NewList[string]()
		//shellScript.Add("source /root/.bashrc")
		shellScript.Add("mkdir -p " + DistRoot + receiver.appGit.GetRelativePath())
		shellScript.Add("cd " + DistRoot + receiver.appGit.GetRelativePath())
		shellScript.Add("set -xe")
		shellScript.AddArray(step.Run)
		shellScript.Add("")
		script := shellScript.ToString("\n")
		// 支持参数化脚本
		for k, v := range step.With {
			script = strings.ReplaceAll(script, "{{"+k+"}}", parse.ToString(v))
		}
		shellPath := fmt.Sprintf("%s%d-%d.sh", ShellRoot, receiver.Env.BuildNumber, index+1)
		file.WriteString(shellPath, script)
		_, wait := receiver.dockerClient.Container.Cp(receiver.fopsBuildName, shellPath, shellPath, receiver.ctx)
		wait()

		// 执行脚本 docker exec FOPS-Build /bin/bash -c "xxx.sh"
		lstResult, wait := receiver.dockerClient.Container.Exec(receiver.fopsBuildName, shellPath, receiver.WorkflowsAction.Env, receiver.ctx)
		code := exec.SaveToChan(receiver.logQueue.progress, lstResult, wait)
		receiver.checkResult(code == 0)
	}
	receiver.logQueue.progress <- "---------------------------------------------------------"
}

// 得到应用的CommitId
func (receiver *BuildEO) getCommitId() {
	cmd := fmt.Sprintf("git -C %s rev-parse HEAD", receiver.appGit.GetAbsolutePath())
	//  docker exec FOPS-Build /bin/bash -c "git -C /var/lib/fops/git/fops rev-parse HEAD"

	progress, wait := receiver.dockerClient.Container.Exec(receiver.fopsBuildName, cmd, nil, receiver.ctx)
	if wait() == 0 {
		if commitId := collections.NewListFromChan(progress).First(); len(commitId) >= 16 {
			receiver.Env.CommitId = commitId[:16]
			receiver.logQueue.progress <- fmt.Sprintf("应用的CommitId：%s", receiver.Env.CommitId)
		}
	}
}

// 得到整个目录的Sha256sum
func (receiver *BuildEO) getSha256sum() {
	// find /var/lib/fops/dist -type f ! -path "*/.git/*" ! -name ".gitignore" ! -name ".gitmodules" ! -name "with.json" -exec sha256sum {} + |sort -k2|sha256sum
	cmd := fmt.Sprintf("find %s -type f ! -path \"*/.git/*\" ! -name \".gitignore\" ! -name \".gitmodules\" ! -name \"with.json\" -exec sha256sum {} + |sort -k2|sha256sum", DistRoot)
	// docker exec FOPS-Build /bin/bash -c "find /var/lib/fops/dist -type f ! -path \"*/.git/*\" ! -name \".gitignore\" ! -name \".gitmodules\" ! -name \"with.json\" -exec sha256sum {} + |sort -k2|sha256sum"
	progress, wait := receiver.dockerClient.Container.Exec(receiver.fopsBuildName, cmd, nil, receiver.ctx)
	if wait() == 0 {
		if sha256sum := collections.NewListFromChan(progress).First(); len(sha256sum) >= 16 {
			receiver.Env.Sha256sum = sha256sum[:16]
			receiver.logQueue.progress <- fmt.Sprintf("打包目录Sha256：%s", receiver.Env.Sha256sum)
		}
	}
}

// 使用缓存
func (receiver *BuildEO) useCache(index int, gits collections.List[GitEO]) bool {
	// 找到该commitId，如果之前存在，则直接使用
	dockerImage := container.Resolve[appsBranch.Repository]().GetDockerImage(receiver.AppName, receiver.Env.Sha256sum)
	if dockerImage == "" {
		return false
	}
	// 如果存在dockerPush、dockerswarmUpdateVer，则直接使用上次构建的镜像
	dockerPushStepVO := collections.NewList(receiver.WorkflowsAction.Steps...).Find(func(item *stepVO) bool {
		return item.ActionName == "dockerPush"
	})
	dockerswarmUpdateVerStepVO := collections.NewList(receiver.WorkflowsAction.Steps...).Find(func(item *stepVO) bool {
		return item.ActionName == "dockerswarmUpdateVer"
	})

	// 如果没有找到dockerPush、dockerswarmUpdateVer，则不使用缓存
	if dockerPushStepVO == nil || dockerswarmUpdateVerStepVO == nil {
		return false
	}

	// 将之前的镜像更新成新的镜像名称
	cmd := fmt.Sprintf("docker tag %s %s", dockerImage, receiver.Env.DockerImage)
	lstResult, wait := receiver.dockerClient.Container.Exec(receiver.fopsBuildName, cmd, nil, receiver.ctx)
	code := exec.SaveToChan(receiver.logQueue.progress, lstResult, wait)
	if code != 0 {
		return false
	}

	receiver.logQueue.progress <- "使用上一次构建的镜像包：" + dockerImage
	receiver.runStep(index+1, *dockerPushStepVO, gits)
	receiver.runStep(index+2, *dockerswarmUpdateVerStepVO, gits)
	return true
}

// GenerateWorkflowsContent 生成Workflows（并更新集群ID）
func (receiver *BuildEO) GenerateWorkflowsContent() bool {
	// 更新工作流文件到本地
	if isSuccess := receiver.gitDevice.PullWorkflows(receiver.ctx, receiver.apps.GetWorkflowsRoot(), receiver.appGit.Branch, receiver.appGit.GetAuthHub(), receiver.logQueue.progress); !isSuccess {
		return false
	}

	// 读取工作流定义的内容
	var err error
	receiver.WorkflowsAction, err = LoadWorkflows(receiver.apps.GetWorkflowsDir()+receiver.WorkflowsName+".yml", receiver.AppName, receiver.appGit.GetName())
	receiver.ClusterId = receiver.WorkflowsAction.ClusterId
	if err != nil {
		receiver.logQueue.progress <- err.Error()
		return false
	}

	// 默认使用本地集群
	if receiver.ClusterId == 0 {
		receiver.ClusterId = container.Resolve[cluster.Repository]().GetLocalCluster().Id
	}

	// 加载Git代理
	if receiver.WorkflowsAction.With["proxy"] != "" {
		receiver.WorkflowsAction.Steps = append([]stepVO{
			{
				Name:              "开启Git代理",
				ActionName:        "gitProxy",
				ActionVer:         "v1",
				ActionDownloadUrl: "https://github.com/farseers/FOPS-Actions/releases/download/v1/gitProxy",
				RepositoryName:    "farseers/FOPS-Actions",
				With:              make(map[string]any),
			},
		}, receiver.WorkflowsAction.Steps...)
	}

	// 加载初始化环境
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

	return true
}

// 将全局参数 覆盖到 系统参数
func (receiver *BuildEO) ReplaceSysWith(sysWith map[string]any) {
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
	for i, step := range receiver.WorkflowsAction.Steps {
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

		// 超时设置
		if t, isOk := step.With["timeout"]; isOk {
			step.Timeout = parse.ToInt(t)
		}
		// 没有设置，则默认5分钟
		if step.Timeout == 0 {
			step.Timeout = 300
		}

		receiver.WorkflowsAction.Steps[i] = step
	}
}

func (receiver *BuildEO) catch() {
	if err := recover(); err != nil {
		if receiver.cancel != nil {
			receiver.cancel()
		}
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
	if status == eumBuildStatus.Cancel {
		exception.ThrowRefuseException("手动取消，退出构建。")
	}

	if !result {
		exception.ThrowRefuseException("exit")
	}
}

// 设置任务失败
func (receiver *BuildEO) fail() {
	// 发布事件
	event.BuildFinishedEvent{AppName: receiver.AppName, BuildId: receiver.Id, ClusterId: receiver.ClusterId, IsSuccess: false}.PublishEvent()

	receiver.logQueue.progress <- "---------------------------------------------------------"
	receiver.logQueue.progress <- "执行失败，退出构建。"

	// 重启构建环境
	receiver.dockerClient.Container.Restart(receiver.fopsBuildName)

	// 更新本次构建状态 = 失败
	container.Resolve[Repository]().SetFail(receiver.Id, receiver.Env)
	queue.Push("monitor", monitor.NewDataEO(receiver.AppName, "build", fmt.Sprintf("分支%s 构建失败", receiver.Env.BranchName)))
}

// 设置任务成功
func (receiver *BuildEO) success() {
	// 包含dockerswarmUpdateVer，才要发布通知
	if collections.NewList(receiver.WorkflowsAction.Steps...).Where(func(item stepVO) bool {
		return item.ActionName == "dockerswarmUpdateVer"
	}).Any() {
		receiver.logQueue.progress <- "更新镜像版本完成。"
		// 发布事件
		event.BuildFinishedEvent{AppName: receiver.AppName, BuildId: receiver.Id, ClusterId: receiver.ClusterId, IsSuccess: true, DockerVer: receiver.Env.BuildNumber, DockerImage: receiver.DockerImage}.PublishEvent()
	}

	receiver.logQueue.progress <- "---------------------------------------------------------"
	receiver.logQueue.progress <- "构建完成。"

	// 更新本次构建状态 = 成功
	receiver.Status = eumBuildStatus.Finish
	receiver.IsSuccess = true
	receiver.FinishAt = dateTime.Now()
	container.Resolve[Repository]().SetSuccess(receiver.Id, receiver.Env)
	queue.Push("monitor", monitor.NewDataEO(receiver.AppName, "build", fmt.Sprintf("分支%s 构建成功：%s", receiver.Env.BranchName, receiver.Env.DockerImage)))
}

// 得到所有Git
func (receiver *BuildEO) getGits() collections.List[GitEO] {
	gits := collections.NewList[GitEO]()
	if !receiver.appGit.IsNil() {
		gits.Add(receiver.appGit)
	}
	// 依赖的框架
	frameworkGits := container.Resolve[Repository]().ToGitList(receiver.apps.FrameworkGits)
	gits.AddList(frameworkGits)
	return gits
}

// GenerateEnv 生成环境变量
func (receiver *BuildEO) GenerateEnv(projectGitRoot string, dockerHub string, dockerImage string, gitName string) {
	receiver.Env = EnvVO{
		BuildId:     receiver.Id,
		BuildNumber: receiver.BuildNumber,
		AppName:     receiver.AppName,
		DockerHub:   dockerHub,
		ClusterId:   receiver.ClusterId,
		DockerImage: dockerImage,
		AppGitRoot:  projectGitRoot,
		GitHub:      receiver.appGit.Hub,
		GitName:     gitName,
		BranchName:  receiver.BranchName,
	}
}

// SetCancel 取消任务
func (receiver *BuildEO) SetCancel() {
	// 更新本次构建状态 = 失败
	container.Resolve[Repository]().SetCancel(receiver.Id, receiver.Env)
}

// WatchStatus 监控当前构建
func (receiver *BuildEO) WatchStatus() {
	for {
		time.Sleep(3 * time.Second)
		curBuildEO := container.Resolve[Repository]().ToBuildEntity(receiver.Id)
		if curBuildEO.Status == eumBuildStatus.Cancel {
			receiver.Status = eumBuildStatus.Cancel
			if receiver.cancel != nil && !curBuildEO.IsSuccess {
				receiver.cancel()
			}
			return
		}
	}
}
