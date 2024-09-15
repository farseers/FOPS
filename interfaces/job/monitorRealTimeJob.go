package job

import (
	"fmt"
	"fops/domain/enum/ruleTimeType"
	"fops/domain/monitor"
	"fops/interfaces/notice"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/tasks"
	"strings"
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
	// 循环map

	for appId, ruleArray := range appIdMap {

		for _, rule := range ruleArray {

			// 获取app数据
			appData := monitorRepository.ToListDataByAppIdKey(appId, rule.KeyName, 1)
			dataInfo := appData.First()
			reqVal := dataInfo.Value

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
				if comparisonMsg != "" && rule.NoticeWhatsAppIds != "" {
					// 通知数据
					noticeList := monitorRepository.ToListNoticeById(strings.Split(rule.NoticeWhatsAppIds, ","))
					noticeList.Foreach(func(not *monitor.NoticeEO) {
						notice.WhatsAppSendMsg(comparisonMsg, not.Phone, not.ApiKey)
					})
				}
			}
			// 监控程序是否正常
			if dateTime.Now().Sub(dataInfo.CreateAt).Minutes() > 10 {
				time.Sleep(1 * time.Second)
				var comparisonMsg = fmt.Sprintf("%s %s %s %s", time.Now().Format("01-02 15:04:05"), rule.AppId, rule.AppName, "请检查项目是否已经停止")
				// 通知数据
				noticeList := monitorRepository.ToListNoticeById(strings.Split(rule.NoticeWhatsAppIds, ","))
				noticeList.Foreach(func(not *monitor.NoticeEO) {
					notice.WhatsAppSendMsg(comparisonMsg, not.Phone, not.ApiKey)
				})
			}
		}
	}

}
