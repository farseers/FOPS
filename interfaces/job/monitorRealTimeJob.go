package job

import (
	"fmt"
	"fops/domain/enum/ruleTimeType"
	"fops/domain/monitor"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/tasks"
	"time"
)

// MonitorRealTimeJob 处理提交过来的监控数据
func MonitorRealTimeJob(*tasks.TaskContext) {
	// 规则
	monitorRepository := container.Resolve[monitor.Repository]()
	// 应用规则数据
	ruleList := monitorRepository.ToListRule()
	appIdMap := make(map[string][]monitor.RuleEO)
	ruleList.GroupBy(&appIdMap, func(item monitor.RuleEO) any {
		return item.AppId
	})
	// 添加通知日志
	addLogs := collections.NewList[monitor.NoticeLogEO]()
	// 循环map
	for appId, ruleArray := range appIdMap {
		var ruleInfo monitor.RuleEO
		for _, rule := range ruleArray {
			ruleInfo = rule
			// 获取app数据
			appData := monitorRepository.ToListDataByAppIdKey(appId, rule.KeyName, 1)
			dataInfo := appData.First()
			reqVal := dataInfo.Value
			ruleInfo = rule
			// 时间类型判断
			var send = false
			switch rule.TimeType {
			case ruleTimeType.Hour:
				if time.Now().Hour() >= rule.StartTime.Hour() && time.Now().Hour() <= rule.EndTime.Hour() {
					send = true
				}
			case ruleTimeType.Day:
				if time.Now().After(rule.StartTime) && time.Now().Before(rule.EndTime) {
					send = true
				}
			}
			if !rule.IsNull() && send {
				var comparisonMsg string
				switch rule.Comparison {
				case ">":
					if parse.ToFloat32(rule.KeyValue) > parse.ToFloat32(reqVal) {
						comparisonMsg = fmt.Sprintf("%s %s %s 大于 %f", time.Now().Format("01-02 15:04:05"), rule.AppId, rule.AppName, parse.ToFloat32(reqVal))
					}
				case "<":
					if parse.ToFloat32(rule.KeyValue) < parse.ToFloat32(reqVal) {
						comparisonMsg = fmt.Sprintf("%s %s %s 小于 %f", time.Now().Format("01-02 15:04:05"), rule.AppId, rule.AppName, parse.ToFloat32(reqVal))
					}
				case "=":
					if parse.ToFloat32(rule.KeyValue) == parse.ToFloat32(reqVal) {
						comparisonMsg = fmt.Sprintf("%s %s %s 等于 %f", time.Now().Format("01-02 15:04:05"), rule.AppId, rule.AppName, parse.ToFloat32(reqVal))
						//comparisonMsg = fmt.Sprintf("%s %s 等于", time.Now().Format("01-02 15:04:05"), rule.AppName)
					}

				}
				// 发送消息 whatsapp
				if comparisonMsg != "" && len(rule.NoticeIds) > 0 {
					// 通知数据
					noticeList := monitorRepository.ToListNoticeById(rule.NoticeIds)
					noticeList.Foreach(func(not *monitor.NoticeEO) {
						not.Notice(comparisonMsg)
						// 记录日志
						addLogs.Add(monitor.NewLog(rule.AppId, rule.AppName, not.Id, not.NoticeType, comparisonMsg))
					})
				}
			}
		}
		// appid 取最大的时间
		maxTime := monitorRepository.GetMaxTimeByAppId(appId)
		// 监控程序是否正常
		if time.Now().Sub(maxTime).Minutes() > 10 {
			time.Sleep(1 * time.Second)
			var comparisonMsg = fmt.Sprintf("%s %s %s %s", time.Now().Format("01-02 15:04:05"), ruleInfo.AppId, ruleInfo.AppName, "请检查项目是否已经停止")
			// 通知数据
			noticeList := monitorRepository.ToListNoticeById(ruleInfo.NoticeIds)
			noticeList.Foreach(func(not *monitor.NoticeEO) {
				not.Notice(comparisonMsg)
				// 记录日志
				addLogs.Add(monitor.NewLog(ruleInfo.AppId, ruleInfo.AppName, not.Id, not.NoticeType, comparisonMsg))
			})
		}
	}

	// 保存日志
	err := monitorRepository.SaveLog(addLogs)
	exception.ThrowWebExceptionError(403, err)
}
