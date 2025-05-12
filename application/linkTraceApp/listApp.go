// @area /linkTrace/
package linkTraceApp

import (
	"fops/domain/linkTrace"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/fs/trace/eumTraceType"
)

// WebApiList WebApi链路追踪列表
// @get webApiList
// @filter application.Jwt
func WebApiList(traceId, appName, appIp, requestIp, searchUrl string, statusCode int, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[trace.TraceContext] {
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

// Delete 删除7天之前的日志
// @post delete
// @filter application.Jwt
func Delete(traceType int, linkTraceRepository linkTrace.Repository) {
	err := linkTraceRepository.Delete(eumTraceType.Enum(traceType), time.Now().AddDate(0, 0, -3))
	exception.ThrowWebExceptionError(403, err)
}

// WebSocketList WebSocket链路追踪列表
// @get webSocketList
// @filter application.Jwt
func WebSocketList(traceId, appName, appIp, requestIp, searchUrl string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[trace.TraceContext] {
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

	return linkTraceRepository.ToWebSocketList(traceId, appName, appIp, requestIp, searchUrl, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}

// TaskList Task链路追踪列表
// @get taskList
// @filter application.Jwt
func TaskList(traceId, appName, appIp, taskName string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[trace.TraceContext] {
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
// @filter application.Jwt
func FScheduleList(traceId, appName, appIp, taskName string, taskGroupId, taskId, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[trace.TraceContext] {
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
// @filter application.Jwt
func ConsumerList(traceId, appName, appIp, server, queueName, routingKey string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[trace.TraceContext] {
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

// QueueList Queue链路追踪列表
// @get queueList
// @filter application.Jwt
func QueueList(traceId, appName, appIp, queueName, routingKey string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[trace.TraceContext] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	queueName = strings.TrimSpace(queueName)
	routingKey = strings.TrimSpace(routingKey)

	return linkTraceRepository.ToQueueList(traceId, appName, appIp, queueName, routingKey, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}

// EventList Event链路追踪列表
// @get eventList
// @filter application.Jwt
func EventList(traceId, appName, appIp, queueName, routingKey string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int, linkTraceRepository linkTrace.Repository) collections.PageList[trace.TraceContext] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	queueName = strings.TrimSpace(queueName)
	routingKey = strings.TrimSpace(routingKey)

	return linkTraceRepository.ToEventList(traceId, appName, appIp, queueName, routingKey, searchUseTs, onlyViewException, startMin, pageSize, pageIndex)
}
