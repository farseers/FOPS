// @area /flog/
package flogApp

import (
	"fops/domain/logData"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"strings"
	"time"
)

// List 日志列表
// @get list
// @filter application.Jwt
func List(traceId, appName, appIp, logContent string, minute int, logLevel eumLogLevel.Enum, pageSize, pageIndex int, logDataRepository logData.Repository) collections.PageList[flog.LogData] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	traceId = strings.TrimSpace(traceId)
	appName = strings.TrimSpace(appName)
	appIp = strings.TrimSpace(appIp)
	logContent = strings.TrimSpace(logContent)

	lst := logDataRepository.ToList(traceId, appName, appIp, logContent, minute, logLevel, pageSize, pageIndex)
	lst.List = lst.List.OrderBy(func(item flog.LogData) any {
		return item.CreateAt.UnixNano()
	}).ToList()
	return lst
}

// Delete 删除7天之前的日志
// @post delete
// @filter application.Jwt
func Delete(logDataRepository logData.Repository) {
	err := logDataRepository.DeleteLog(time.Now().AddDate(0, 0, -7))
	exception.ThrowWebExceptionError(403, err)
}

// Info 日志详情
// @get info-{id}
// @filter application.Jwt
func Info(id string, logDataRepository logData.Repository) flog.LogData {
	return logDataRepository.ToInfo(id)
}

// StatCount 日志类型统计
// @get StatCount
// @filter application.Jwt
func StatCount(appName string, logDataRepository logData.Repository) collections.List[logData.LogCountEO] {
	appName = strings.TrimSpace(appName)
	return logDataRepository.StatCount()
}
