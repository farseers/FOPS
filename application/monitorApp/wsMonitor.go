// @area /ws/
package monitorApp

import (
	"fops/domain/monitor"
	"time"

	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snc"
	fsMonitor "github.com/farseer-go/monitor"
	"github.com/farseer-go/queue"
	"github.com/farseer-go/webapi/websocket"
)

// WsReceive 监控数据接收
// @ws monitor
func WsReceive(context *websocket.Context[fsMonitor.SendContentVO], monitorRepository monitor.Repository) {
	// 如果appId为空直接返回
	context.ForReceiverFunc(func(req *fsMonitor.SendContentVO) {
		if len(req.AppId) == 0 || len(req.AppName) == 0 {
			return
		}
		jsonData, _ := snc.Marshal(req)
		flog.Info("WsReceive:" + string(jsonData))
		// 所有key值进行处理
		req.Keys.Keys().Foreach(func(key *string) {
			reqVal := req.Keys.GetValue(*key)
			// 添加消息队列
			queue.Push("monitor", monitor.NewDataEO(req.AppName, *key, parse.ToString(reqVal)))
		})
	})
}

// WsNotice 通知消息
// @ws notice
func WsNotice(context *websocket.Context[string], monitorRepository monitor.Repository) {

	context.ReceiverMessageFunc(5*time.Second, func(message string) {
		// 未读消息列表
		noReadList := ToListPageNoticeLogNoRead(20, monitorRepository)
		sendMap := make(map[string]interface{})
		sendMap["NoReadList"] = noReadList
		err := context.Send(sendMap)
		exception.ThrowWebExceptionError(403, err)
	})

}
