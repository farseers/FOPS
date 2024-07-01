package localQueue

import (
	"fops/domain/logData"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
)

func SaveFlogQueue(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	lst := collections.NewList[flog.LogData]()
	lstMessage.Foreach(func(item *any) {
		lst.Add((*item).(flog.LogData))
	})

	err := container.Resolve[logData.Repository]().Save(lst)
	flog.ErrorIfExists(err)
}
