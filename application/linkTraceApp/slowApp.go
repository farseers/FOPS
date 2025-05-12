// @area /linkTrace/
package linkTraceApp

import (
	"fops/domain/linkTrace"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
)

// DeleteSlow 删除7天之前的日志
// @post deleteSlow
// @filter application.Jwt
func DeleteSlow(dbName string, linkTraceRepository linkTrace.Repository) {
	err := linkTraceRepository.DeleteSlow(dbName, time.Now().AddDate(0, 0, -3))
	exception.ThrowWebExceptionError(403, err)
}

// SlowDbList 慢数据库列表
// @get slowDbList
// @filter application.Jwt
func SlowDbList(traceId, appName, appIp, dbName, tableName string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[linkTrace.TraceDetailEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	dbName = strings.TrimSpace(dbName)
	tableName = strings.TrimSpace(tableName)

	return linkTraceRepository.ToSlowDbList(traceId, appName, appIp, dbName, tableName, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}

// SlowEsList 慢Es列表
// @get slowEsList
// @filter application.Jwt
func SlowEsList(traceId, appName, appIp, indexName, aliasesName string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[linkTrace.TraceDetailEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	indexName = strings.TrimSpace(indexName)
	aliasesName = strings.TrimSpace(aliasesName)

	return linkTraceRepository.ToSlowEsList(traceId, appName, appIp, indexName, aliasesName, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}

// SlowEtcdList 慢Etcd列表
// @get slowEtcdList
// @filter application.Jwt
func SlowEtcdList(traceId, appName, appIp, key string, leaseID int64, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[linkTrace.TraceDetailEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	key = strings.TrimSpace(key)

	return linkTraceRepository.ToSlowEtcdList(traceId, appName, appIp, key, leaseID, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}

// SlowHandList 慢手动列表
// @get slowHandList
// @filter application.Jwt
func SlowHandList(traceId, appName, appIp, name string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[linkTrace.TraceDetailEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	name = strings.TrimSpace(name)

	return linkTraceRepository.ToSlowHandList(traceId, appName, appIp, name, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}

// SlowHttpList 慢Http列表
// @get slowHttpList
// @filter application.Jwt
func SlowHttpList(traceId, appName, appIp, method, url, body string, statusCode int, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[linkTrace.TraceDetailEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	method = strings.TrimSpace(method)
	url = strings.TrimSpace(url)
	body = strings.TrimSpace(body)

	return linkTraceRepository.ToSlowHttpList(traceId, appName, appIp, method, url, body, statusCode, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}

// SlowMqList 慢Mq列表
// @get slowMqList
// @filter application.Jwt
func SlowMqList(traceId, appName, appIp, server, exchange, routingKey string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[linkTrace.TraceDetailEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	server = strings.TrimSpace(server)
	exchange = strings.TrimSpace(exchange)
	routingKey = strings.TrimSpace(routingKey)

	return linkTraceRepository.ToSlowMqList(traceId, appName, appIp, server, exchange, routingKey, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}

// SlowRedisList 慢Redis列表
// @get slowRedisList
// @filter application.Jwt
func SlowRedisList(traceId, appName, appIp, methodName, key, field string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[linkTrace.TraceDetailEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	methodName = strings.TrimSpace(methodName)
	key = strings.TrimSpace(key)
	field = strings.TrimSpace(field)

	return linkTraceRepository.ToSlowRedisList(traceId, appName, appIp, methodName, key, field, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}
