package job

import (
	"fops/domain/apps"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/tasks"
)

// CollectsClusterJob 3秒收集一次Docker集群信息
func CollectsClusterJob(*tasks.TaskContext) {
	dockerSwarmDevice := container.Resolve[apps.IDockerSwarmDevice]()
	appsRepository := container.Resolve[apps.Repository]()

	container.Resolve[core.ITransaction]("default").Transaction(func() {
		// 收集所有节点的信息
		nodeList := dockerSwarmDevice.NodeList()
		nodeList.Foreach(func(node *apps.DockerNodeVO) {
			vo := dockerSwarmDevice.NodeInfo(node.NodeName)
			node.IP = vo.IP
			node.OS = vo.OS
			node.Architecture = vo.Architecture
			node.CPUs = vo.CPUs
			node.Memory = vo.Memory
			node.Label = vo.Label
		})
		appsRepository.UpdateClusterNode(nodeList)

		// 收集所有服务的运行情况
		serviceList := dockerSwarmDevice.ServiceList()
		_, _ = appsRepository.UpdateInsReplicas(serviceList)
	})
}
