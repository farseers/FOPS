// @area /apps/
package appsApp

import (
	"fops/application/appsApp/response"
	"fops/domain/appsBranch"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
)

// AllBranchList 获取所有应用的分支列表
// @post autobuild/list
// @filter application.Jwt
func AllBranchList(appsBranchRepository appsBranch.Repository) collections.List[response.BranchListResponse] {
	var mGroupBy map[string]collections.List[appsBranch.DomainObject]
	appsBranchRepository.ToListByAutoBuild().GroupBy(&mGroupBy, func(item appsBranch.DomainObject) any {
		return item.AppName
	})

	lst := collections.NewList[response.BranchListResponse]()
	for appName, list := range mGroupBy {
		lst.Add(response.BranchListResponse{
			AppName: appName,
			List:    list,
		})
	}

	return lst.OrderBy(func(item response.BranchListResponse) any {
		return item.AppName
	}).ToList()
}

// BuildList 获取指定应用的分支列表
// @post autobuild/branchList
// @filter application.Jwt
func BranchList(appName string, appsBranchRepository appsBranch.Repository) collections.List[appsBranch.DomainObject] {
	return appsBranchRepository.ToListByAppName(appName)
}

// ResetCommitId 重置错误次数
// @post autobuild/resetCommitId
// @filter application.Jwt
func ResetCommitId(commitId string, appsBranchRepository appsBranch.Repository) {
	err := appsBranchRepository.ResetCommitId(commitId)
	exception.ThrowWebExceptionError(403, err)
}
