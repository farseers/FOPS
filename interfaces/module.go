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
	}

	tasks.Run("统计应用的在线实例", time.Second*20, job.StatAppOnlineJob, context.Background())
	tasks.RunNow("清除历史注册信息", time.Hour*1, job.ClearHistoryRegisterAppJob, context.Background())
}
