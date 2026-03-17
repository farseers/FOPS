package repository

import (
	"fops/domain/apps"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
)

type appsFrameworkRepository struct {
}

// ToAppsFrameworkList 获取应用的框架列表
func (receiver *appsFrameworkRepository) ToAppsFrameworkList(appName string) collections.List[apps.AppsFrameworkEO] {
	lst := context.MysqlContext.AppsFramework.Where("app_name = ?", appName).ToList()
	return mapper.ToList[apps.AppsFrameworkEO](lst)
}

// AddAppsFramework 添加应用框架关系
func (receiver *appsFrameworkRepository) AddAppsFramework(eo apps.AppsFrameworkEO) error {
	po := mapper.Single[model.AppsFrameworkPO](eo)
	return context.MysqlContext.AppsFramework.Insert(&po)
}

// UpdateAppsFramework 更新应用框架关系
func (receiver *appsFrameworkRepository) UpdateAppsFramework(eo apps.AppsFrameworkEO) (int64, error) {
	po := mapper.Single[model.AppsFrameworkPO](eo)
	return context.MysqlContext.AppsFramework.Where("app_name = ? and framework_id = ?", eo.AppName, eo.FrameworkId).Update(po)
}

// DeleteAppsFramework 删除应用框架关系
func (receiver *appsFrameworkRepository) DeleteAppsFramework(appName string, frameworkId int64) (int64, error) {
	return context.MysqlContext.AppsFramework.Where("app_name = ? and framework_id = ?", appName, frameworkId).Delete()
}

// DeleteAppsFrameworkByAppName 删除应用的所有框架关系
func (receiver *appsFrameworkRepository) DeleteAppsFrameworkByAppName(appName string) (int64, error) {
	return context.MysqlContext.AppsFramework.Where("app_name = ?", appName).Delete()
}

// UpdateCommitId 更新框架的CommitId
func (receiver *appsFrameworkRepository) UpdateCommitId(appName string, frameworkId int64, commitId string) (int64, error) {
	return context.MysqlContext.AppsFramework.Where("app_name = ? and framework_id = ?", appName, frameworkId).UpdateValue("commit_id", commitId)
}

// ExistsAppsFramework 判断应用框架关系是否存在
func (receiver *appsFrameworkRepository) ExistsAppsFramework(appName string, frameworkId int64) bool {
	return context.MysqlContext.AppsFramework.Where("app_name = ? and framework_id = ?", appName, frameworkId).IsExists()
}
