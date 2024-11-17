package job

import (
	"fops/domain/clusterNode"
	"time"

	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
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
		dockerNode := dockerClient.Node.Info(node.NodeName)
		node.IsHealth = node.Status == "Ready" && node.Availability == "Active"
		node.IP = dockerNode.IP
		node.OS = dockerNode.OS
		node.Architecture = dockerNode.Architecture
		node.CPUs = dockerNode.CPUs
		node.Memory = dockerNode.Memory
		node.Label = dockerNode.Label

		if node.IP == "" {
			flog.Warningf("集群节点：%s，没有读取到IP", node.NodeName)
			return
		}

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

	// 存在没有收集到IP的情况时，退出
	if dockerNodeList.Where(func(item docker.DockerNodeVO) bool {
		return item.IP == ""
	}).Any() {
		return
	}

	// 删除旧的节点
	clusterNode.NodeList.Foreach(func(dockerNodeVO *docker.DockerNodeVO) {
		// 如果不在docker swarm中了，说明机器从集群中删除了。
		if !dockerNodeList.Where(func(dockerItem docker.DockerNodeVO) bool {
			return dockerItem.IP == dockerNodeVO.IP
		}).Any() {
			flog.Warningf("集群节点：%s 已离开集群", dockerNodeVO.IP)
			clusterNode.NodeList.RemoveAll(func(item docker.DockerNodeVO) bool {
				return item.IP == dockerNodeVO.IP
			})
		}
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
