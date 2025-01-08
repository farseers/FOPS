package repository

import (
	"fops/domain/appsBranch"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"

	"github.com/farseer-go/data"
	"github.com/farseer-go/mapper"
)

type appsBranchRepository struct {
	// IRepository 通用的仓储接口
	data.IRepository[appsBranch.DomainObject]
}

// UpdateByBranch 更新
func (receiver *appsBranchRepository) UpdateByBranch(do appsBranch.DomainObject) error {
	po := mapper.Single[model.AppsBranchPO](do)
	_, err := context.MysqlContext.AppsBranch.Where("app_name = ? and branch_name = ?", do.AppName, do.BranchName).Update(po)
	return err
}

// DeleteBranch 删除分支
func (receiver *appsBranchRepository) DeleteBranch(appName, branchName string) error {
	_, err := context.MysqlContext.AppsBranch.Where("app_name = ? and branch_name = ?", appName, branchName).Delete()
	return err
}

// GetUnRunUT 获取未运行UT的分支
func (receiver *appsBranchRepository) GetUnRunUT() appsBranch.DomainObject {
	po := context.MysqlContext.AppsBranch.Where("build_success = 0 and build_error_count < 3").Asc("commit_at").ToEntity()
	return mapper.Single[appsBranch.DomainObject](po)
}
