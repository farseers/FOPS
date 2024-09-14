// @area /ws/
package monitorApp

import (
	"fops/domain/monitor"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	fsMonitor "github.com/farseer-go/monitor"
	"github.com/farseer-go/webapi/websocket"
)

// WsReceive 监控数据接收
// @ws monitor
func WsReceive(context *websocket.Context[fsMonitor.SendContentVO], monitorRepository monitor.Repository) {
	req := context.Receiver()
	flog.Info(req)
	// 如果appId为空直接返回
	if len(req.AppId) == 0 {
		return
	}
	// 规则
	ruleList := monitorRepository.ToListRuleByAppId(req.AppId)
	// 所有key值进行处理
	// 添加记录
	addList := collections.NewList[monitor.DataEO]()
	req.Keys.Keys().Foreach(func(key *string) {
		ruleVal := ruleList.Where(func(rule monitor.RuleEO) bool {
			return rule.KeyName == *key
		}).First()
		reqVal := req.Keys.GetValue(*key)
		// 发送消息
		if len(ruleVal.NoticeWhatsAppApiKey) > 0 {
			ruleVal.WhatsAppSendMsg(parse.ToString(reqVal))
		}

		addList.Add(monitor.NewDataEO(req.AppId, req.AppName, *key, parse.ToString(reqVal)))
	})
	err := monitorRepository.Save(addList)
	exception.ThrowWebExceptionError(403, err)
}
