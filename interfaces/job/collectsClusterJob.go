package job

import (
	"fops/domain/apps"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/tasks"
)

// CollectsClusterJob 3秒收集一次Docker集群信息
func CollectsClusterJob(*tasks.TaskContext) {
	dockerSwarmDevice := container.Resolve[apps.IDockerSwarmDevice]()
	serviceList := dockerSwarmDevice.ServiceList()
	serviceList.Foreach(func(item *apps.DockerName) {
		flog.Printf("%+v", item)
	})
}
