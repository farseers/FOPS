package configure

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
)

// Repository 仓储接口
type Repository interface {
	// IRepository 通用的仓储接口
	data.IRepository[DomainObject]
	// ToListByAppName 获取应用的配置列表
	ToListByAppName(appName string) collections.List[DomainObject]
	// Rollback 回滚版本
	Rollback(appName string) (int64, error)
}
