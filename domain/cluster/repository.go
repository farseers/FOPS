package cluster

import (
	"github.com/farseer-go/data"
)

// Repository 仓储接口
type Repository interface {
	// IRepository 通用的仓储接口
	data.IRepository[DomainObject]
	CancelLocal(id int64)          // 设置其它集群为非本地
	GetLocalCluster() DomainObject // 获取本地集群
}
