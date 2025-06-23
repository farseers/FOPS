package job

import (
	"fops/domain/clusterNode"
	"time"

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
			return
		}

		node.IP = dockerNode.IP
		node.Label = dockerNode.Label
		// 以下属性暂时不需要，由代理节点获取
		// node.OS = dockerNode.OS
		// node.Architecture = dockerNode.Architecture
		// node.CPUs = dockerNode.CPUs
		// node.Memory = dockerNode.Memory

		// 加入到本地列表
		dockerNodeVO := clusterNode.NodeList.Find(func(dockerItem *docker.DockerNodeVO) bool {
			return dockerItem.IP == node.IP
		})
		if dockerNodeVO == nil {
			flog.Infof("发现新的集群节点：%s", node.IP)
			clusterNode.NodeList.Add(*node)
		} else {
			dockerNodeVO.IsHealth = node.IsHealth
			dockerNodeVO.Label = node.Label
			dockerNodeVO.UpdateAt = time.Now()
			dockerNodeVO.Status = node.Status
			dockerNodeVO.Availability = node.Availability
			// 以下属性暂时不需要，由代理节点获取
			// dockerNodeVO.OS = node.OS
			// dockerNodeVO.Architecture = node.Architecture
			// dockerNodeVO.CPUs = node.CPUs
			// dockerNodeVO.Memory = node.Memory
		}
	})
}
