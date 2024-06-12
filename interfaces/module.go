package interfaces

import (
	"context"
	"fops/application"
	"fops/domain/apps"
	"fops/interfaces/job"
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
	dockerVer := container.Resolve[apps.IDockerDevice]().GetVersion()
	if dockerVer != "" {
		tasks.Run("开启构建", time.Second*1, job.BuildingJob, context.Background())
		flog.Info("Docker version：" + flog.Blue(dockerVer))

		// 3秒收集一次Docker集群信息
		tasks.Run("收集Docker集群信息", time.Second*3, job.CollectsClusterJob, context.Background())
	}

	tasks.RunNow("统计访问", time.Minute*1, job.StatVisitsJob, context.Background())
}
