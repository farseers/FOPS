// @area /linkTrace/
package linkTraceApp

import (
	"fops/domain/linkTrace"
	"github.com/farseer-go/collections"
	linkTraceCom "github.com/farseer-go/linkTrace"
	"strings"
)

// WebApiList WebApi链路追踪列表
// @get webApiList
func WebApiList(traceId, appName, appIp, requestIp, searchUrl string, statusCode int, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[linkTraceCom.TraceContext] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	requestIp = strings.TrimSpace(requestIp)
	searchUrl = strings.TrimSpace(searchUrl)

	return linkTraceRepository.ToWebApiList(traceId, appName, appIp, requestIp, searchUrl, statusCode, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}

// TaskList Task链路追踪列表
// @get taskList
func TaskList(traceId, appName, appIp, taskName string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[linkTraceCom.TraceContext] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	taskName = strings.TrimSpace(taskName)

	return linkTraceRepository.ToTaskList(traceId, appName, appIp, taskName, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}

// FScheduleList FSchedule链路追踪列表
// @get fScheduleList
func FScheduleList(traceId, appName, appIp, taskName string, taskGroupId, taskId, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[linkTraceCom.TraceContext] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	taskName = strings.TrimSpace(taskName)

	return linkTraceRepository.ToFScheduleList(traceId, appName, appIp, taskName, taskGroupId, taskId, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}

// ConsumerList Consumer链路追踪列表
// @get consumerList
func ConsumerList(traceId, appName, appIp, server, queueName, routingKey string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[linkTraceCom.TraceContext] {
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
	queueName = strings.TrimSpace(queueName)
	routingKey = strings.TrimSpace(routingKey)

	return linkTraceRepository.ToConsumerList(traceId, appName, appIp, server, queueName, routingKey, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}
