package localQueue

import (
	"fops/domain/enum/ruleTimeType"
	"fops/domain/monitor"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/trace"
)

func SaveMonitorDataQueue(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	lstData := collections.NewList[monitor.DataEO]()
	for _, item := range lstMessage.ToArray() {
		data := item.(*monitor.DataEO)
		lstData.Add(*data)
	}

	// 规则
	monitorRepository := container.Resolve[monitor.Repository]()
	// 应用规则数据
	ruleList := monitorRepository.ToListRule()
	// 添加通知日志
	addLogs := collections.NewList[monitor.NoticeLogEO]()
	addDataLogs := collections.NewList[monitor.DataEO]()
	appNameList := collections.NewList[string]()
	lstData.Foreach(func(dataEO *monitor.DataEO) {
		//flog.Info("消息队列处理")
		//flog.Info(dataEO)
		if !appNameList.Contains(dataEO.AppName) {
			appNameList.Add(dataEO.AppName)
		}
		rules := ruleList.Where(func(item monitor.RuleEO) bool {
			return strings.Contains(item.AppName, dataEO.AppName) && item.KeyName == dataEO.Key
		}).ToList()
		// 是否保存数据
		var isSaveData = false
		// 规则列表
		rules.Foreach(func(rule *monitor.RuleEO) {
			// 应用名称集合
			reqVal := dataEO.Value
			// 时间类型判断
			var send = false
			switch rule.TimeType {
			case ruleTimeType.Hour:
				startTime := parse.ToInt(strings.ReplaceAll(rule.StartDate, ":", ""))
				endTime := parse.ToInt(strings.ReplaceAll(rule.EndDate, ":", ""))
				nowTime := parse.ToInt(time.Now().Format("150405"))
				if nowTime >= startTime && nowTime <= endTime {
					send = true
					isSaveData = true
				}
			case ruleTimeType.Day:
				startDay := parse.ToTime(rule.StartDay)
				endDay := parse.ToTime(rule.EndDay)
				if time.Now().After(startDay) && time.Now().Before(endDay) {
					send = true
					isSaveData = true
				}
			}
			if !rule.IsNull() && send {
				comparisonMsg := ""
				// 比较结果
				if rule.CompareResult(reqVal) {
					comparisonMsg = rule.GetTipTemplate(dataEO.AppName, reqVal)
				}
				// 发送消息 whatsapp
				if len(comparisonMsg) > 0 && len(rule.NoticeIds) > 0 {
					// 通知数据
					noticeList := monitorRepository.ToListNoticeById(rule.NoticeIds)
					noticeList.Foreach(func(not *monitor.NoticeEO) {
						not.Notice(comparisonMsg)
						// 记录日志
						addLogs.Add(monitor.NewLog(dataEO.AppName, not.Id, not.Name, not.NoticeType, comparisonMsg))
					})
				}
			}
		})
		// 记录数据
		if isSaveData {
			addDataLogs.Add(*dataEO)
		}
	})

	// 保存日志
	if addLogs.Count() > 0 {
		err := monitorRepository.SaveLog(addLogs)
		exception.ThrowWebExceptionError(403, err)
	}
	// 保存数据
	if addDataLogs.Count() > 0 {
		err := monitorRepository.Save(addDataLogs)
		exception.ThrowWebExceptionError(403, err)
	}
	// 刷新时间
	if appNameList.Count() > 0 {
		appNameList.Foreach(func(item *string) {
			if !monitorRepository.IsExistSyncAt(*item) {
				err := monitorRepository.SaveSyncAt(monitor.NewSyncAtEO(*item))
				exception.ThrowWebExceptionError(403, err)
			} else {
				// 更新时间
				err := monitorRepository.UpdateSyncAt(*item, time.Now())
				exception.ThrowWebExceptionError(403, err)
			}
		})
	}
}
