// @area /ws/
package monitorApp

import (
	"fmt"
	"fops/domain/monitor"
	"fops/interfaces/notice"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	fsMonitor "github.com/farseer-go/monitor"
	"github.com/farseer-go/webapi/websocket"
	"strings"
	"time"
)

// WsReceive 监控数据接收
// @ws monitor
func WsReceive(context *websocket.Context[fsMonitor.SendContentVO], monitorRepository monitor.Repository) {
	req := context.Receiver()
	// 如果appId为空直接返回
	if len(req.AppId) == 0 {
		return
	}
	// 规则
	ruleList := monitorRepository.ToListRuleByAppId(req.AppId)
	// 所有key值进行处理
	req.Keys.Keys().Foreach(func(key *string) {
		ruleVal := ruleList.Where(func(rule monitor.RuleEO) bool {
			return rule.KeyName == *key
		}).First()
		reqVal := req.Keys.GetValue(*key)
		if !ruleVal.IsNull() {
			var comparisonMsg string
			switch ruleVal.Comparison {
			case ">":
				if parse.ToFloat32(ruleVal.KeyValue) > parse.ToFloat32(reqVal) {
					comparisonMsg = "大于"
				}
			case "<":
				if parse.ToFloat32(ruleVal.KeyValue) < parse.ToFloat32(reqVal) {
					comparisonMsg = "小于"
				}
			case "=":
				if parse.ToFloat32(ruleVal.KeyValue) == parse.ToFloat32(reqVal) {
					comparisonMsg = "等于"
				}
			}
			// 发送消息
			if comparisonMsg != "" && ruleVal.NoticeWhatsAppApiKey != "" {
				notice.WhatsAppSendMsg(fmt.Sprintf("%s %s：%s，%s：%s", time.Now().Format("2006-01-02 15:04:05"), ruleVal.AppName, ruleVal.Remark, comparisonMsg, reqVal), strings.Split(ruleVal.NoticeWhatsAppApiKey, ","))
			}
		}
	})

	// 添加记录
	err := monitorRepository.Add(monitor.NewLogEO(req.AppId, req.AppName, req.Keys))
	exception.ThrowWebExceptionError(403, err)
}
