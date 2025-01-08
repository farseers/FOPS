package clusterNode

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
)

// Repository 集群节点
type Repository interface {
	// IRepository 通用的仓储接口
	UpdateClusterNode(lst collections.List[docker.DockerNodeVO])                                                                                                   // 更新集群节点信息
	GetClusterNodeList() collections.List[docker.DockerNodeVO]                                                                                                     // 获取集群节点列表
	UpdateClusterNodeResourceByAgentIP(agentIP string, cpuUsagePercent, memoryUsagePercent, memoryUsage float64, disk uint64, diskUsagePercent, diskUsage float64) // 更新集群节点的资源信息
	Delete(ip string)                                                                                                                                              // 删除集群节点
}
