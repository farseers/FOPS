package consumer

import (
	"encoding/json"
	"fops/domain/enum/ruleTimeType"
	"fops/domain/monitor"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/rabbit"
	"strings"
	"time"
)

// MonitorDataConsumer 监控数据处理
func MonitorDataConsumer(messages collections.List[rabbit.EventArgs]) bool {
	// 监控数据
	lstData := collections.NewList[monitor.DataEO]()
	messages.Foreach(func(item *rabbit.EventArgs) {
		var dataEO monitor.DataEO
		_ = json.Unmarshal(item.Body, &dataEO)
		lstData.Add(dataEO)
	})
	// 规则
	monitorRepository := container.Resolve[monitor.Repository]()
	// 应用规则数据
	ruleList := monitorRepository.ToListRule()
	// 添加通知日志
	addLogs := collections.NewList[monitor.NoticeLogEO]()
	addDataLogs := collections.NewList[monitor.DataEO]()
	appNameList := collections.NewList[string]()
	lstData.Foreach(func(dataEO *monitor.DataEO) {
		if !appNameList.Contains(dataEO.AppName) {
			appNameList.Add(dataEO.AppName)
		}
		rules := ruleList.Where(func(item monitor.RuleEO) bool {
			return strings.Contains(item.AppName, dataEO.AppName)
		}).ToList()
		// 规则列表
		rules.Foreach(func(rule *monitor.RuleEO) {
			// 应用名称集合
			reqVal := dataEO.Value
			// 时间类型判断
			var send = false
			switch rule.TimeType {
			case ruleTimeType.Hour:
				startTime := parse.ToInt(rule.StartTime.Format("150405"))
				endTime := parse.ToInt(rule.EndTime.Format("150405"))
				nowTime := parse.ToInt(time.Now().Format("150405"))
				if nowTime >= startTime && nowTime <= endTime {
					send = true
				}
			case ruleTimeType.Day:
				if time.Now().After(rule.StartTime) && time.Now().Before(rule.EndTime) {
					send = true
				}
			}
			if !rule.IsNull() && send {
				comparisonMsg := ""
				switch rule.Comparison {
				case ">":
					if parse.ToFloat32(reqVal) > parse.ToFloat32(rule.KeyValue) {
						comparisonMsg = rule.GetTipTemplate(dataEO.AppName, reqVal)
					}
				case "<":
					if parse.ToFloat32(reqVal) < parse.ToFloat32(rule.KeyValue) {
						comparisonMsg = rule.GetTipTemplate(dataEO.AppName, reqVal)
					}
				case "=":
					if parse.ToFloat32(rule.KeyValue) == parse.ToFloat32(reqVal) {
						comparisonMsg = rule.GetTipTemplate(dataEO.AppName, reqVal)
					}
				}
				// 发送消息 whatsapp
				if len(comparisonMsg) > 0 && len(rule.NoticeIds) > 0 {
					// 通知数据
					noticeList := monitorRepository.ToListNoticeById(rule.NoticeIds)
					noticeList.Foreach(func(not *monitor.NoticeEO) {
						not.Notice(comparisonMsg)
						// 记录日志
						addLogs.Add(monitor.NewLog(rule.AppName, not.Id, not.Name, not.NoticeType, comparisonMsg))
						// 记录数据
						addDataLogs.Add(*dataEO)
					})
				}
			}
		})
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
			}
		})
	}
	return true
}
