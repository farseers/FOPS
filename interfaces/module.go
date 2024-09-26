package interfaces

import (
	"context"
	"fops/application"
	"fops/domain/_/eumBuildStatus"
	"fops/domain/apps"
	"fops/domain/apps/event"
	"fops/interfaces/job"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/tasks"
	"github.com/farseer-go/webapi"
	"time"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{webapi.Module{}, application.Module{}}
}

func (module Module) PostInitialize() {
	client := docker.NewClient()
	dockerVer := client.GetVersion()
	if dockerVer != "" {
		tasks.Run("开启构建", time.Second*1, job.BuildingJob, context.Background())
		flog.Info("Docker version：" + flog.Blue(dockerVer))

		// 3秒收集一次Docker集群信息
		tasks.Run("收集Docker集群信息", time.Second*3, job.CollectsClusterJob, context.Background())
	}

	tasks.Run("统计访问", time.Minute*1, job.StatVisitsJob, context.Background())

	tasks.Run("监控数据处理", time.Second*5, job.MonitorRealTimeJob, context.Background())
	tasks.Run("fops监控数据处理", time.Second*1, job.MonitorFopsJob, context.Background())

	// 监听agent的IP变化
	go job.ListenerAgentNotify()

	// 如果最后一次构建是fops，且状态=构建中，同时fops的仓库=最后一次构建的镜像，则强制做一次同步操作
	buildEO := container.Resolve[apps.Repository]().GetLastBuilding()
	appEO := container.Resolve[apps.Repository]().ToEntity("fops")
	if buildEO.AppName == appEO.AppName && buildEO.Status == eumBuildStatus.Building && appEO.DockerImage == buildEO.DockerImage {
		// 发布事件
		event.BuildFinishedEvent{AppName: appEO.AppName, BuildId: buildEO.Id, ClusterId: buildEO.ClusterId, IsSuccess: true}.PublishEvent()
		container.Resolve[apps.Repository]().SetSuccessForFops(buildEO.Id)
	}
}
