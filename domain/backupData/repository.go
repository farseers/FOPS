package backupData

import (
	"github.com/farseer-go/data"
)

// Repository 仓储接口
type Repository interface {
	// IRepository 通用的仓储接口
	data.IRepository[DomainObject]
	GetCountById(id string) int64 // 查看同名的ID的数量
}
