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
	// 所有key值进行处理
	// 添加记录
	addList := collections.NewList[monitor.DataEO]()
	req.Keys.Keys().Foreach(func(key *string) {
		reqVal := req.Keys.GetValue(*key)
		addList.Add(monitor.NewDataEO(req.AppName, *key, parse.ToString(reqVal)))
	})
	err := monitorRepository.Save(addList)
	exception.ThrowWebExceptionError(403, err)
}