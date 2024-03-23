package infrastructure

import (
	"fops/domain/apps/event"
	"fops/domain/linkTrace"
	"fops/infrastructure/device"
	"fops/infrastructure/domainEvent"
	"fops/infrastructure/localQueue"
	"fops/infrastructure/repository"
	"fops/infrastructure/repository/context"
	"github.com/farseer-go/data"
	"github.com/farseer-go/data/driver/clickhouse"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/modules"
	linkTraceModule "github.com/farseer-go/linkTrace"
	"github.com/farseer-go/queue"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	// 这些模块都是farseer-go内置的模块
	return []modules.FarseerModule{data.Module{}, eventBus.Module{}, linkTraceModule.Module{}}
}

func (module Module) PostInitialize() {
	// 初始化数据库上下文
	context.InitMysqlContext()
	// 初始化仓储
	repository.InitRepository()

	// 注册驱动
	device.RegisterDockerDevice()
	device.RegisterDockerSwarmDevice()
	device.RegisterKubectlDevice()

	eventBus.RegisterEvent(event.BuildFinishedEventName).RegisterSubscribe("更新应用的版本信息", domainEvent.BuildFinishedConsumer)
	eventBus.RegisterEvent(event.DockerPushedEventName).RegisterSubscribe("docker推送完成事件", domainEvent.DockerPushedConsumer)
	eventBus.RegisterEvent(event.GitCloneOrPulledEventName).RegisterSubscribe("更新git拉取时间", domainEvent.GitCloneOrPulledConsumer)

	// 启用链路追踪写入CH
	linkTrace.Config = configure.ParseConfig[linkTrace.ConfigEO]("Fops.LinkTrace")
	if linkTrace.Config.Driver == "clickhouse" {
		data_clickhouse.Module{}.Initialize()
		context.InitChContextContext()
	}

	// 日志消费
	queue.Subscribe("flog", "saveToCh", 1000, localQueue.SaveFlogQueue)
}
