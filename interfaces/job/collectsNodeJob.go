package job

import (
	"fops/domain/clusterNode"

	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
)

func CollectsNodeJob(*tasks.TaskContext) {
	//appsRepository := container.Resolve[apps.Repository]()
	//clusterRepository := container.Resolve[cluster.Repository]()
	clusterNodeRepository := container.Resolve[clusterNode.Repository]()

	// 收集所有节点的信息
	dockerClient := docker.NewClient()
	dockerNodeList := dockerClient.Node.List()
	// 没有读取到集群，则退出
	if dockerNodeList.Count() == 0 {
		return
	}
	dockerNodeList.Foreach(func(node *docker.DockerNodeVO) {
		vo := dockerClient.Node.Info(node.NodeName)
		node.IsHealth = node.Status == "Ready" && node.Availability == "Active"
		node.IP = vo.IP
		node.OS = vo.OS
		node.Architecture = vo.Architecture
		node.CPUs = vo.CPUs
		node.Memory = vo.Memory
		node.Label = vo.Label
	})

	// 数据库的节点列表
	clusterNodeList := clusterNodeRepository.GetClusterNodeList()
	// 删除旧的节点
	clusterNodeList.Foreach(func(dockerNodeVO *docker.DockerNodeVO) {
		if !dockerNodeList.Where(func(dockerItem docker.DockerNodeVO) bool {
			return dockerItem.IP == dockerNodeVO.IP
		}).Any() {
			clusterNodeRepository.Delete(dockerNodeVO.IP)
			clusterNodeList.RemoveAll(func(item docker.DockerNodeVO) bool {
				return item.IP == dockerNodeVO.IP
			})
		}
	})

	// 将新的节点加入节点列表

	// 读取数据库的节点列表，删除离开的节点
	clusterNodeRepository.GetClusterNodeList().Foreach(func(item *docker.DockerNodeVO) {
		if !dockerNodeList.Where(func(dockerItem docker.DockerNodeVO) bool {
			return dockerItem.IP == item.IP
		}).Any() {
			clusterNodeRepository.Delete(item.IP)
		}
	})

}