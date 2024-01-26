package apps

import (
	"context"
	"fops/domain/_/eumBuildStatus"
	"fops/domain/apps/event"
	"fops/domain/cluster"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/utils/file"
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
	BuildServerId     int64               // 构建的服务端id
	Log               []string            // 构建日志
	Env               EnvVO               // 环境变量
	AppName           string              // 应用名称
	Dockerfile        string              // Dockerfile内容
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

	// 打印环境变量
	receiver.Env.Print(receiver.logQueue.progress)

	// 前置检查
	receiver.directoryDevice.Check(receiver.logQueue.progress)

	// 拉取主仓库及依赖仓库
	receiver.checkResult(receiver.gitDevice.CloneOrPullAndDependent(receiver.getGits(), receiver.logQueue.progress, receiver.ctx))

	// 登陆镜像仓库(先登陆，如果失败了，后则面也不需要编译、打包了)
	receiver.checkResult(receiver.dockerDevice.Login(clusterDO.DockerHub, clusterDO.DockerUserName, clusterDO.DockerUserPwd, receiver.logQueue.progress, receiver.Env, receiver.ctx))

	// 将需要打包的源代码，复制到dist目录
	receiver.copyToDistDevice.Copy(receiver.getGits(), receiver.Env, receiver.logQueue.progress)

	// 生成Dockerfile文件
	receiver.checkResult(receiver.GenerateDockerfileContent())
	receiver.dockerDevice.CreateDockerfile(receiver.AppName, receiver.Dockerfile, receiver.ctx)

	// docker打包
	receiver.checkResult(receiver.dockerDevice.Build(receiver.Env, receiver.logQueue.progress, receiver.ctx))

	// docker上传
	receiver.checkResult(receiver.dockerDevice.Push(receiver.Env, receiver.logQueue.progress, receiver.ctx))

	// 首次创建还是更新镜像
	if receiver.dockerSwarmDevice.ExistsDocker(clusterDO, receiver.AppName) {
		// 更新镜像
		//receiver.checkResult(receiver.kubectlDevice.SetImages(receiver.Cluster, receiver.AppName, receiver.Env.DockerImage, receiver.Project.K8SControllersType, receiver.progress, receiver.ctx))
		receiver.checkResult(receiver.dockerSwarmDevice.SetImages(clusterDO, receiver.AppName, receiver.Env.DockerImage, receiver.logQueue.progress, receiver.ctx))
	} else {
		// 创建容器服务
		receiver.checkResult(receiver.dockerSwarmDevice.CreateService(receiver.AppName, receiver.apps.DockerNodeRole, receiver.apps.AdditionalScripts, clusterDO.DockerNetwork, receiver.apps.DockerReplicas, receiver.Env.DockerImage, receiver.logQueue.progress, receiver.ctx))
	}
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

// GenerateDockerfileContent 替换模板
func (receiver *BuildEO) GenerateDockerfileContent() bool {
	// 为空时，读取应用git文件中的Dockerfile文件
	if receiver.Dockerfile == "" {
		// 如果没有自定义，则使用应用仓库根目录的Dockerfile文件
		if receiver.apps.DockerfilePath == "" {
			receiver.apps.DockerfilePath = receiver.appGit.GetAbsolutePath() + "Dockerfile"
		} else {
			// 自定义Dockerfile路径
			if strings.HasPrefix(receiver.apps.DockerfilePath, "/") {
				receiver.apps.DockerfilePath = receiver.apps.DockerfilePath[1:]
			} else if strings.HasPrefix(receiver.apps.DockerfilePath, "./") {
				receiver.apps.DockerfilePath = receiver.apps.DockerfilePath[2:]
			}
			receiver.apps.DockerfilePath = receiver.appGit.GetAbsolutePath() + receiver.apps.DockerfilePath
		}
		receiver.Dockerfile = file.ReadString(receiver.apps.DockerfilePath)
		if receiver.Dockerfile == "" {
			receiver.logQueue.progress <- "Dockerfile没有定义。"
			return false
		}
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
