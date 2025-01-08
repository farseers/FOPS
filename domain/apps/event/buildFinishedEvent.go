package event

import "github.com/farseer-go/eventBus"

// BuildFinishedEventName 事件名称
const BuildFinishedEventName = "BuildFinishedEvent"

// BuildFinishedEvent 构建完成后，发布事件
type BuildFinishedEvent struct {
	// 项目ID
	AppName string
	// 本次构建ID
	BuildId int64
	// 构建的集群
	ClusterId int64
	// 是否成功
	IsSuccess bool
}

// PublishEvent 发布事件
func (receiver BuildFinishedEvent) PublishEvent() {
	_ = eventBus.PublishEvent(BuildFinishedEventName, receiver)
}
