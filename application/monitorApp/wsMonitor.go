// @area /ws/
package monitorApp

import (
	"encoding/json"
	"fops/domain/monitor"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
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
		jsonData, _ := json.Marshal(req)
		flog.Info("WsReceive:" + string(jsonData))
		// 所有key值进行处理
		req.Keys.Keys().Foreach(func(key *string) {
			reqVal := req.Keys.GetValue(*key)
			// 添加消息队列
			queue.Push("monitor", monitor.NewDataEO(req.AppName, *key, parse.ToString(reqVal)))
		})
	})
}
