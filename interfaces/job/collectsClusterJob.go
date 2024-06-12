package job

import (
	"fmt"
	"fops/domain/apps"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
	"strings"
)

// CollectsClusterJob 3秒收集一次Docker集群信息
func CollectsClusterJob(*tasks.TaskContext) {
	dockerSwarmDevice := container.Resolve[apps.IDockerSwarmDevice]()
	serviceList := dockerSwarmDevice.ServiceList()
	serviceList.Foreach(func(service *string) {
		sers := strings.Split(*service, "\t")
		fmt.Println(sers[0])
	})
}
