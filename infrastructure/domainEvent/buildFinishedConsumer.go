package domainEvent

import (
	"fops/domain/apps"
	"fops/domain/apps/event"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
)

// BuildFinishedConsumer 构建完成事件：更新应用的版本信息
func BuildFinishedConsumer(message any, _ core.EventArgs) {
	appsRepository := container.Resolve[apps.Repository]()
	buildFinishedEvent := message.(event.BuildFinishedEvent)

	flog.Infof("收到构建完成事件消息：IsSuccess=%v,ClusterId=%d,BuildId=%d", buildFinishedEvent.IsSuccess, buildFinishedEvent.ClusterId, buildFinishedEvent.BuildId)

	// 更新项目的版本信息
	appsDO := appsRepository.ToEntity(buildFinishedEvent.AppName)
	appsDO.UpdateBuildVer(buildFinishedEvent.IsSuccess, buildFinishedEvent.ClusterId, buildFinishedEvent.BuildId)

	flog.Infof("更新后的集群map=%+v", appsDO.ClusterVer)
	_, _ = appsRepository.UpdateClusterVer(buildFinishedEvent.AppName, appsDO.ClusterVer)
}
