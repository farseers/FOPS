package backupData

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
)

// Repository 仓储接口
type Repository interface {
	// IRepository 通用的仓储接口
	data.IRepository[DomainObject]
	GetCountById(id string) int64                       // 查看同名的ID的数量
	ToNextBackupData() DomainObject                     // 获取即将备份的数据
	AddHistory(lst collections.List[BackupHistoryData]) // 添加备份文件列表
	UpdateAt(id any, do DomainObject) (int64, error)    // 更新备份计划的时间
}
