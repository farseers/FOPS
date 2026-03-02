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
		dockerNode := dockerClient.Node.Info(node.Description.Hostname)
		if dockerNode.Status.Addr == "" {
			return
		}

		// 加入到本地列表
		dockerNodeVO := clusterNode.NodeList.Find(func(dockerItem *docker.DockerNodeVO) bool {
			return dockerItem.Status.Addr == node.Status.Addr
		})
		if dockerNodeVO == nil {
			flog.Infof("发现新的集群节点：%s", node.Status.Addr)
			clusterNode.NodeList.Add(*node)
			// 重新排序
			clusterNode.NodeList = clusterNode.NodeList.OrderBy(func(item docker.DockerNodeVO) any {
				return item.Status.Addr
			}).ToList()
		} else {
			dockerNodeVO.IsHealth = node.IsHealth
			dockerNodeVO.Label = node.Label
			dockerNodeVO.UpdatedAt = time.Now()
			dockerNodeVO.Status = node.Status
			dockerNodeVO.Spec.Availability = node.Spec.Availability
		}
	})
}
