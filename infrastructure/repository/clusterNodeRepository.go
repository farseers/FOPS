package repository

import (
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/mapper"
)

type clusterNodeRepository struct {
}

// UpdateClusterNode 更新集群节点信息
func (receiver *clusterNodeRepository) UpdateClusterNode(lst collections.List[docker.DockerNodeVO]) {
	lstPO := mapper.ToList[model.ClusterNodePO](lst)
	context.MysqlContext.ClusterNode.Where("1=1").Delete()
	context.MysqlContext.ClusterNode.InsertIgnoreList(lstPO, 100)
}

func (receiver *clusterNodeRepository) GetClusterNodeList() collections.List[docker.DockerNodeVO] {
	lstPO := context.MysqlContext.ClusterNode.Desc("is_master").ToList()
	return mapper.ToList[docker.DockerNodeVO](lstPO)
}

func (receiver *clusterNodeRepository) UpdateClusterNodeResourceByAgentIP(agentIP string, cpuUsagePercent, memoryUsagePercent, memoryUsage float64, disk uint64, diskUsagePercent, diskUsage float64) {
	_, _ = context.MysqlContext.ClusterNode.Where("agent_ip = ?", agentIP).Select("cpu_usage_percent", "memory_usage_percent", "memory_usage", "disk", "disk_usage_percent", "disk_usage").Update(model.ClusterNodePO{
		CpuUsagePercent:    cpuUsagePercent,
		MemoryUsagePercent: memoryUsagePercent,
		MemoryUsage:        memoryUsage,
		Disk:               disk,
		DiskUsagePercent:   diskUsagePercent,
		DiskUsage:          diskUsage,
	})
}
func (receiver *clusterNodeRepository) Delete(ip string) {
	_, _ = context.MysqlContext.ClusterNode.Where("ip = ?", ip).Delete()
}
