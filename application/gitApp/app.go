// @area /git/
package gitApp

import (
	"fops/application/gitApp/request"
	"fops/domain/apps"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
	"strings"
)

// Add 添加Git
// @post add
// @filter application.Jwt
func Add(req request.AddRequest, appsRepository apps.Repository) {
	do := mapper.Single[apps.GitEO](req)
	if !strings.HasPrefix(do.Path, "/") {
		do.Path = "/" + do.Path
	}
	// 添加
	err := appsRepository.AddGit(do)
	exception.ThrowWebExceptionError(403, err)
}

// Update 修改Git
// @post update
// @filter application.Jwt
func Update(req request.UpdateRequest, appsRepository apps.Repository) {
	do := mapper.Single[apps.GitEO](req)
	exception.ThrowWebExceptionBool(!appsRepository.ExistsGit(req.Id), 403, "GitId不存在")

	if !strings.HasPrefix(do.Path, "/") {
		do.Path = "/" + do.Path
	}
	if !strings.HasSuffix(do.Path, "/") {
		do.Path += "/"
	}
	_, err := appsRepository.UpdateGit(do)
	exception.ThrowWebExceptionError(403, err)
}

// List Git列表
// @post list
// @filter application.Jwt
func List(isApp int, appsRepository apps.Repository) collections.List[apps.GitEO] {
	return appsRepository.ToGitListAll(isApp)
}

// Delete 删除Git
// @post delete
// @filter application.Jwt
func Delete(gitId int64, appsRepository apps.Repository) {
	exception.ThrowWebExceptionBool(gitId < 1, 403, "GitId没有填")
	_, err := appsRepository.DeleteGit(gitId)
	exception.ThrowWebExceptionError(403, err)
}

// Info 查询Git
// @post info
// @filter application.Jwt
func Info(gitId int64, appsRepository apps.Repository) apps.GitEO {
	return appsRepository.ToGitEntity(gitId)
}
