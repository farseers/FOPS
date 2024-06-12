package job

import (
	"fmt"
	"fops/domain/apps"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
	"strings"
)

// CollectsClusterJob 3秒收集一次Docker集群信息
func CollectsClusterJob(*tasks.TaskContext) {
	dockerSwarmDevice := container.Resolve[apps.IDockerSwarmDevice]()
	serviceList := dockerSwarmDevice.ServiceList()
	// 移除标题
	serviceList.RemoveAt(0)
	serviceList.Foreach(func(service *string) {
		// 移除容器ID
		*service = strings.TrimSpace((*service)[12:])
		*service = strings.Replace(*service, "\t", "", -1)
		sers := collections.NewList(strings.Split(*service, " ")...)
		sers.RemoveAll(func(item string) bool {
			return item == ""
		})
		fmt.Println(sers.ToString("|"))
	})
}
