package job

import (
	"fops/domain/clusterNode"
	"time"

	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
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

	// 删除旧的节点
	clusterNode.NodeList.Foreach(func(dockerNodeVO *docker.DockerNodeVO) {
		dockerNode := dockerNodeList.Find(func(dockerItem *docker.DockerNodeVO) bool {
			return dockerItem.IP == dockerNodeVO.IP
		})

		// 如果不在docker swarm中了，说明机器从集群中删除了。
		if dockerNode == nil {
			clusterNode.NodeList.RemoveAll(func(item docker.DockerNodeVO) bool {
				return item.IP == dockerNodeVO.IP
			})
			return
		}

		// 更新状态
		dockerNodeVO.IsHealth = dockerNode.IsHealth
		dockerNodeVO.OS = dockerNode.OS
		dockerNodeVO.Architecture = dockerNode.Architecture
		dockerNodeVO.CPUs = dockerNode.CPUs
		dockerNodeVO.Memory = dockerNode.Memory
		dockerNodeVO.Label = dockerNode.Label
	})

	// 间隔更新
	if time.Now().Second()%5 == 0 {
		// 通过事务来更新
		container.Resolve[core.ITransaction]("default").Transaction(func() {
			// 更新集群节点信息
			clusterNodeRepository.UpdateClusterNode(clusterNode.NodeList)
		})
	}
}
