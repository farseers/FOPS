package interfaces

import (
	"context"
	"fops/application"
	"fops/domain/_/eumBuildStatus"
	"fops/domain/_/eumBuildType"
	"fops/domain/apps"
	"fops/domain/apps/event"
	"fops/interfaces/job"
	"time"

	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/tasks"
	"github.com/farseer-go/webapi"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{webapi.Module{}, application.Module{}}
}

func (module Module) PostInitialize() {
	client := docker.NewClient()
	if client.IsMaster() {
		tasks.Run("开启构建应用", time.Second*1, job.BuildAppJob, context.Background())
		tasks.Run("开启自动构建", time.Second*1, job.AutoBuildAppJob, context.Background())
		flog.Info("Docker version：" + flog.Blue(client.GetVersion()))

		// 3秒收集一次Docker集群信息
		tasks.RunNow("收集Docker Swarm集群信息", time.Second*10, job.CollectsNodeJob, context.Background())
		tasks.Run("收集Docker应用信息", time.Second*3, job.CollectsClusterJob, context.Background())

		tasks.Run("统计访问", time.Minute*1, job.StatVisitsJob, context.Background())

		tasks.Run("fops监控数据处理", time.Minute*1, job.MonitorFopsJob, context.Background())
		tasks.Run("同步Git分支", time.Second*30, job.SyncAppsBranchJob, context.Background())
		tasks.Run("执行备份计划", time.Minute*1, job.SyncBackupDataJob, context.Background())

		// 监听agent的IP变化
		go job.ListenerAgentNotify()

		// 如果最后一次构建是fops，且状态=构建中，同时fops的仓库=最后一次构建的镜像，则强制做一次同步操作
		buildEO := container.Resolve[apps.Repository]().GetLastBuilding(eumBuildType.Manual)
		appEO := container.Resolve[apps.Repository]().ToEntity("fops")
		if buildEO.AppName == appEO.AppName && buildEO.Status == eumBuildStatus.Building && appEO.DockerImage == buildEO.DockerImage {
			flog.Infof("恭喜，你正在使用最新的FOPS版本：%s", appEO.DockerImage)
			// 发布事件
			event.BuildFinishedEvent{AppName: appEO.AppName, BuildId: buildEO.Id, ClusterId: buildEO.ClusterId, IsSuccess: true}.PublishEvent()
			container.Resolve[apps.Repository]().SetSuccessForFops(buildEO.Id)
		}
	}
}
