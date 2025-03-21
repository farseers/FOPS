package job

import (
	"fops/domain/clusterNode"

	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/tasks"
)

// 收集Docker swarm集群节点信息
func CollectsNodeJob(*tasks.TaskContext) {
	// 收集所有节点的信息
	dockerClient := docker.NewClient()
	dockerNodeList := dockerClient.Node.List()

	// 没有读取到集群，则退出
	if dockerNodeList.Count() == 0 {
		return
	}

	dockerNodeList.Foreach(func(node *docker.DockerNodeVO) {
		// 通过ps节点信息，拿到节点的资源信息
		dockerNode := dockerClient.Node.Info(node.NodeName)
		if dockerNode.IP == "" {
			flog.Warningf("集群节点：%s，没有读取到IP", node.NodeName)
			return
		}

		node.IP = dockerNode.IP
		node.OS = dockerNode.OS
		node.Architecture = dockerNode.Architecture
		node.CPUs = dockerNode.CPUs
		node.Memory = dockerNode.Memory
		node.Label = dockerNode.Label

		// 加入到本地列表
		dockerNodeVO := clusterNode.NodeList.Find(func(dockerItem *docker.DockerNodeVO) bool {
			return dockerItem.IP == node.IP
		})
		if dockerNodeVO == nil {
			flog.Infof("发现新的集群节点：%s", node.IP)
			clusterNode.NodeList.Add(*node)
		} else {
			dockerNodeVO.IsHealth = node.IsHealth
			dockerNodeVO.OS = node.OS
			dockerNodeVO.Architecture = node.Architecture
			dockerNodeVO.CPUs = node.CPUs
			dockerNodeVO.Memory = node.Memory
			dockerNodeVO.Label = node.Label
		}
	})
}
