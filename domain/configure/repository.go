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
	// ToEntityByKey 获取对象
	ToEntityByKey(appName, key string) DomainObject
	// GetLastVer 获取最后一个版本号
	GetLastVer(appName, key string) int
	// Rollback 回滚版本
	Rollback(appName, key string, ver int) (int64, error)
	// DeleteKey 删除Key
	DeleteKey(appName, key string) (int64, error)
}
