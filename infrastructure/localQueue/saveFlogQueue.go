package localQueue

import (
	"fops/domain/logData"
	"fops/domain/monitor"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/queue"
)

func SaveFlogQueue(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}
	lst := collections.NewList[flog.LogData]()
	for _, item := range lstMessage.ToArray() {
		data := item.(*flog.LogData)
		lst.Add(*data)

		switch data.LogLevel {
		case eumLogLevel.Error:
			queue.Push("monitor", monitor.NewDataEO(data.AppName, "log_error", data.Content))
		case eumLogLevel.Warning:
			queue.Push("monitor", monitor.NewDataEO(data.AppName, "log_warning", data.Content))
		}
	}

	err := container.Resolve[logData.Repository]().Save(lst)
	flog.ErrorIfExists(err)
}
