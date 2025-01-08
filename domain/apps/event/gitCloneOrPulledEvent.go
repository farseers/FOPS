package event

import "github.com/farseer-go/eventBus"

const GitCloneOrPulledEventName = "GitCloneOrPulledEvent"

type GitCloneOrPulledEvent struct {
	// GitId
	GitId int
}

// PublishEvent 发布事件
func (receiver GitCloneOrPulledEvent) PublishEvent() {
	_ = eventBus.PublishEvent(GitCloneOrPulledEventName, receiver)
}
