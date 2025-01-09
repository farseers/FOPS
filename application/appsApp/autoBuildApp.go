// @area /apps/
package appsApp

import (
	"fops/domain/appsBranch"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
)

// BuildList 获取分支列表
// @post autobuild/list
// @filter application.Jwt
func BranchList(appsBranchRepository appsBranch.Repository) map[string]collections.List[appsBranch.DomainObject] {
	var mGroupBy map[string]collections.List[appsBranch.DomainObject]
	appsBranchRepository.ToList().GroupBy(&mGroupBy, func(item appsBranch.DomainObject) any {
		return item.AppName
	})
	return mGroupBy
}

// ResetCommitId 重置错误次数
// @post autobuild/resetCommitId
// @filter application.Jwt
func ResetCommitId(commitId string, appsBranchRepository appsBranch.Repository) {
	err := appsBranchRepository.ResetCommitId(commitId)
	exception.ThrowWebExceptionError(403, err)
}
