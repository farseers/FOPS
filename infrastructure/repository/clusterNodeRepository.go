package repository

import (
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
)

type clusterNodeRepository struct {
}

// UpdateClusterNode 更新集群节点信息
func (receiver *clusterNodeRepository) UpdateClusterNode(lst collections.List[docker.DockerNodeVO]) {
	lstPO := mapper.ToList[model.ClusterNodePO](lst)
	lstPO.Foreach(func(item *model.ClusterNodePO) {
		item.UpdateAt = dateTime.Now()
		// 更新数据
		count, err := context.MysqlContext.ClusterNode.Where("node_name", item.NodeName).Omit("cpu_usage_percent", "memory_usage_percent", "memory_usage", "disk", "disk_usage_percent", "disk_usage").Update(*item)
		flog.ErrorIfExists(err)

		// 没有更新到数据时，则插入
		if count == 0 {
			err = context.MysqlContext.ClusterNode.InsertIgnore(item)
			flog.ErrorIfExists(err)
		}
	})
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
