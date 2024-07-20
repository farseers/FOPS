package repository

import (
	"fops/domain/cluster"
	"fops/infrastructure/repository/context"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/mapper"
)

type clusterRepository struct {
	// IRepository 通用的仓储接口
	data.IRepository[cluster.DomainObject]
}

// CancelLocal 设置其它集群为非本地
func (repository *clusterRepository) CancelLocal(id int64) {
	_, _ = context.MysqlContext.Cluster.Where("id <> ?", id).UpdateValue("is_local", false)
}

// GetLocalCluster 获取本地集群
func (repository *clusterRepository) GetLocalCluster() cluster.DomainObject {
	po := context.MysqlContext.Cluster.Where("is_local = ?", true).ToEntity()
	return mapper.Single[cluster.DomainObject](po)
}

// ToList 获取集群列表
func (repository *clusterRepository) ToList() collections.List[cluster.DomainObject] {
	lstPO := context.MysqlContext.Cluster.Order("is_local desc, id asc").ToList()
	return mapper.ToList[cluster.DomainObject](lstPO)
}
