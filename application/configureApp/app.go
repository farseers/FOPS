// @area /configure/
package configureApp

import (
	"fops/application/configureApp/request"
	"fops/domain/configure"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
)

// List 获取应用的配置列表
// @post list
func List(appName string, configureRepository configure.Repository) collections.List[configure.DomainObject] {
	return configureRepository.ToListByAppName(appName)
}

// AllList 获取所有应用的配置列表
// @get allList
// @filter application.Jwt
func AllList(configureRepository configure.Repository) collections.List[configure.DomainObject] {
	return configureRepository.ToList()
}

// Add 添加配置
// @post add
// @filter application.Jwt
func Add(req request.AddRequest, configureRepository configure.Repository) {
	do := mapper.Single[configure.DomainObject](req)
	do.Ver = 1

	// 添加
	err := configureRepository.Add(do)
	exception.ThrowWebExceptionError(403, err)
}

// Update 修改配置
// @post update
// @filter application.Jwt
func Update(req request.UpdateRequest, configureRepository configure.Repository) {
	do := mapper.Single[configure.DomainObject](req)
	oldDO := configureRepository.ToEntityByKey(do.AppName, do.Key)
	exception.ThrowRefuseExceptionBool(oldDO.IsNil(), "配置不存在")

	// 值相等，不用保存
	if req.Value == oldDO.Value {
		return
	}

	var newDO = configure.DomainObject{
		AppName: req.AppName,
		Key:     req.Key,
		Ver:     oldDO.Ver + 1,
		Value:   req.Value,
	}
	err := configureRepository.Add(newDO)
	exception.ThrowWebExceptionError(403, err)
}

// Rollback 回滚配置
// @post rollback
// @filter application.Jwt
func Rollback(appName, key string, configureRepository configure.Repository) {
	lastVer := configureRepository.GetLastVer(appName, key)
	exception.ThrowWebExceptionfBool(lastVer < 2, 403, "没有可回滚的版本")

	_, err := configureRepository.Rollback(appName, key, lastVer)
	exception.ThrowWebExceptionError(403, err)
}

// Delete 删除配置
// @post delete
// @filter application.Jwt
func Delete(appName, key string, configureRepository configure.Repository) {
	exception.ThrowWebExceptionBool(appName == "", 403, "应用名称没有填")
	exception.ThrowWebExceptionBool(key == "", 403, "Key没有填")
	_, err := configureRepository.DeleteKey(appName, key)
	exception.ThrowWebExceptionError(403, err)
}
