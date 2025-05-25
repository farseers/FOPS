package domainEvent

import (
	"fops/domain/apps"
	"fops/domain/apps/event"

	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
)

// BuildFinishedConsumer 构建完成事件：更新应用的版本信息
func BuildFinishedConsumer(message any, _ core.EventArgs) {
	appsRepository := container.Resolve[apps.Repository]()
	buildFinishedEvent := message.(event.BuildFinishedEvent)

	// 更新项目的版本信息
	appsDO := appsRepository.ToEntity(buildFinishedEvent.AppName)
	appsDO.UpdateBuildVer(buildFinishedEvent.IsSuccess, buildFinishedEvent.ClusterId, buildFinishedEvent.BuildId, buildFinishedEvent.DockerVer, buildFinishedEvent.DockerImage)

	_, _ = appsRepository.UpdateClusterVer(buildFinishedEvent.AppName, appsDO.ClusterVer)
}
