package appsBranch

import "github.com/farseer-go/data"

// Repository 仓储接口
type Repository interface {
	// IRepository 通用的仓储接口
	data.IRepository[DomainObject]
	// UpdateByBranch 更新
	UpdateByBranch(do DomainObject) error
	// 删除分支
	DeleteBranch(appName, branchName string) error
	// GetUnRunUT 获取未运行UT的分支
	GetUnRunUT() DomainObject
}
