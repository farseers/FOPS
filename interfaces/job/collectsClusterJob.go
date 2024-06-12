package job

import (
	"fmt"
	"fops/domain/apps"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
)

// CollectsClusterJob 3秒收集一次Docker集群信息
func CollectsClusterJob(*tasks.TaskContext) {
	dockerSwarmDevice := container.Resolve[apps.IDockerSwarmDevice]()
	serviceList := dockerSwarmDevice.ServiceList()
	if serviceList.Count() == 0 {
		fmt.Println("没有获取到任何服务")
	}
	fmt.Println(serviceList.ToString("\r\n"))
}
