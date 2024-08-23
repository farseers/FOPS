// @area /cluster/
package clusterApp

import (
	"fops/application/clusterApp/request"
	"fops/domain/cluster"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
)

// Add 添加集群
// @post add
// @filter application.Jwt
func Add(req request.AddRequest, clusterRepository cluster.Repository) {
	do := mapper.Single[cluster.DomainObject](req)

	// 如果新添加的集群为本地集群，则取消之前的本地集群
	if req.IsLocal {
		clusterRepository.CancelLocal(0)
	}
	// 添加
	err := clusterRepository.Add(do)
	exception.ThrowWebExceptionError(403, err)
}

// Update 修改集群
// @post update
// @filter application.Jwt
func Update(req request.UpdateRequest, clusterRepository cluster.Repository) {
	do := mapper.Single[cluster.DomainObject](req)
	exception.ThrowWebExceptionBool(!clusterRepository.IsExists(req.Id), 403, "集群Id不存在")

	// 如果新添加的集群为本地集群，则取消之前的本地集群
	if req.IsLocal {
		clusterRepository.CancelLocal(do.Id)
	}
	_, err := clusterRepository.Update(do.Id, do)
	exception.ThrowWebExceptionError(403, err)
}

// List 集群列表
// @post list
// @filter application.Jwt
func List(clusterRepository cluster.Repository) collections.List[cluster.DomainObject] {
	return clusterRepository.ToList()
}

// Delete 删除集群
// @post delete
// @filter application.Jwt
func Delete(clusterId int64, clusterRepository cluster.Repository) {
	exception.ThrowWebExceptionBool(clusterId < 1, 403, "集群Id没有填")
	_, err := clusterRepository.Delete(clusterId)
	exception.ThrowWebExceptionError(403, err)
}
