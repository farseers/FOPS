package accountLogin

import "github.com/farseer-go/data"

// Repository 仓储接口
type Repository interface {
	// IRepository 通用的仓储接口
	data.IRepository[DomainObject]
	// 根据帐号获取数据
	ToEntityByAccountName(loginName string) DomainObject
}
