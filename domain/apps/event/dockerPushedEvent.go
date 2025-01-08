package event

import "github.com/farseer-go/eventBus"

const DockerPushedEventName = "DockerPushedEvent"

type DockerPushedEvent struct {
	// 构建版本号
	BuildNumber int
	// 应用名称
	AppName string
	// 镜像名称
	ImageName string
}

// PublishEvent 发布事件
func (receiver DockerPushedEvent) PublishEvent() {
	_ = eventBus.PublishEvent(DockerPushedEventName, receiver)
}
