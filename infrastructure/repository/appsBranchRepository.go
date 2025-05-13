package repository

import (
	"fops/domain/appsBranch"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/mapper"
)

type appsBranchRepository struct {
	// IRepository 通用的仓储接口
	data.IRepository[appsBranch.DomainObject]
}

// ToListByAppName 获取所有分支
func (receiver *appsBranchRepository) ToListByAutoBuild() collections.List[appsBranch.DomainObject] {
	lstPO := context.MysqlContext.AppsBranch.Where("app_name in (select app_name from apps where ut_workflows_name <>'')").ToList()
	return mapper.ToList[appsBranch.DomainObject](lstPO)
}

// ToListByAppName 获取当前应用的所有分支
func (receiver *appsBranchRepository) ToListByAppName(appName string) collections.List[appsBranch.DomainObject] {
	lstPO := context.MysqlContext.AppsBranch.Where("app_name = ?", appName).Select("branch_name", "commit_at").Desc("commit_at").ToList()
	return mapper.ToList[appsBranch.DomainObject](lstPO)
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
	po := context.MysqlContext.AppsBranch.Where("build_success = 0 and build_error_count < 3 and app_name in (select app_name from apps where ut_workflows_name <>'')").Asc("commit_at").ToEntity()
	return mapper.Single[appsBranch.DomainObject](po)
}

// 重置构建错误
func (receiver *appsBranchRepository) ResetCommitId(commitId string) error {
	_, err := context.MysqlContext.AppsBranch.Select("build_success", "build_error_count", "commit_id", "build_id", "build_at").Where("commit_id = ?", commitId).Update(model.AppsBranchPO{
		BuildSuccess:    false,
		BuildErrorCount: 0,
		CommitId:        commitId,
		BuildId:         0,
		BuildAt:         dateTime.Now(),
	})
	return err
}
