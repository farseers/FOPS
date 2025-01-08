package infrastructure

import (
	"fops/domain/apps/event"
	"fops/domain/fSchedule"
	"fops/domain/linkTrace"
	"fops/infrastructure/device"
	"fops/infrastructure/domainEvent"
	"fops/infrastructure/http"
	"fops/infrastructure/localQueue"
	"fops/infrastructure/repository"
	"fops/infrastructure/repository/context"
	"time"

	"github.com/farseer-go/data"
	data_clickhouse "github.com/farseer-go/data/driver/clickhouse"
	"github.com/farseer-go/eventBus"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/queue"
)

type Module struct {
}

func (module Module) DependsModule() []modules.FarseerModule {
	// 这些模块都是farseer-go内置的模块
	return []modules.FarseerModule{data.Module{}, eventBus.Module{}, queue.Module{}}
}

func (module Module) PostInitialize() {
	container.Register(func() fSchedule.Http { return &http.FScheduleHttp{} })

	// 初始化数据库上下文
	context.InitMysqlContext()
	// 初始化仓储
	repository.InitRepository()

	// 注册驱动
	device.RegisterDockerDevice()
	device.RegisterKubectlDevice()
	device.RegisterGitDevice()

	eventBus.RegisterEvent(event.BuildFinishedEventName).RegisterSubscribe("更新应用的版本信息", domainEvent.BuildFinishedConsumer)
	eventBus.RegisterEvent(event.DockerPushedEventName).RegisterSubscribe("docker推送完成事件", domainEvent.DockerPushedConsumer)
	eventBus.RegisterEvent(event.GitCloneOrPulledEventName).RegisterSubscribe("更新git拉取时间", domainEvent.GitCloneOrPulledConsumer)

	// 启用链路追踪写入CH
	linkTrace.Config = configure.ParseConfig[linkTrace.ConfigEO]("Fops.LinkTrace")
	if linkTrace.Config.Driver == "clickhouse" {
		data_clickhouse.Module{}.Initialize()
		context.InitChContextContext()
	}
	// 监控数据消息处理
	queue.Subscribe("monitor", "saveMonitorData", 1000, 5*time.Second, localQueue.SaveMonitorDataQueue)

	// 日志消费
	queue.Subscribe("flog", "saveFlogDataToCh", 1000, 3*time.Second, localQueue.SaveFlogQueue)

	// 链路追踪日志消费
	queue.Subscribe("linkTrace", "saveLinkTraceLogToCh", 1000, 1*time.Second, localQueue.SaveLinkTraceQueue)
}
