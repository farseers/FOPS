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

	// 收集所有服务的运行情况
	serviceList := dockerSwarmDevice.ServiceList()
	// 如果服务不存在，则添加到列表中，用于更新到数据库中，指明服务的实例为0
	lstApp := appsRepository.ToList()
	lstApp.Foreach(func(appDO *apps.DomainObject) {
		if !serviceList.Where(func(service apps.DockerServiceVO) bool {
			return appDO.AppName == service.Name
		}).Any() {
			serviceList.Add(apps.DockerServiceVO{
				Name:      appDO.AppName,
				Instances: 0,
				Image:     "",
				Replicas:  appDO.DockerReplicas,
			})
		}
	})

	container.Resolve[core.ITransaction]("default").Transaction(func() {
		// 更新集群节点信息
		appsRepository.UpdateClusterNode(nodeList)
		// 更新服务运行情况
		if serviceList.Count() > 0 {
			_, _ = appsRepository.UpdateInsReplicas(serviceList)
		}
	})
}
