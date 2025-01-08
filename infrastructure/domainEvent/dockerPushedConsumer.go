package domainEvent

import (
	"fops/domain/apps"
	"fops/domain/apps/event"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
)

// DockerPushedConsumer docker推送完成事件
func DockerPushedConsumer(message any, ea core.EventArgs) {
	appsRepository := container.Resolve[apps.Repository]()
	dockerPushedEvent := message.(event.DockerPushedEvent)

	// 更新项目的版本信息
	_, _ = appsRepository.UpdateDockerVer(dockerPushedEvent.AppName, dockerPushedEvent.BuildNumber, dockerPushedEvent.ImageName)
}
