package appsBranch

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
)

// Repository 仓储接口
type Repository interface {
	// IRepository 通用的仓储接口
	data.IRepository[DomainObject]
	// ToListByAppName 获取所有分支
	ToListByAutoBuild() collections.List[DomainObject]
	// ToListByAppName 获取当前应用的所有分支
	ToListByAppName(appName string) collections.List[DomainObject]
	// UpdateByBranch 更新
	UpdateByBranch(do DomainObject) error
	// 重置构建错误
	ResetCommitId(commitId string) error
	// 删除分支
	DeleteBranch(appName, branchName string) error
	// GetUnRunUT 获取未运行UT的分支
	GetUnRunUT() DomainObject
	// UpdateDockerImage 更新镜像
	UpdateDockerImage(appName, commitId, dockerImage, sha256sum string) error
	// GetDockerImage 通过sha256sum获取Docker镜像
	GetDockerImage(appName, sha256sum string) string
}
