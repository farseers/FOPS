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
	for _, item := range lstMessage.ToArray() {
		data := item.(*flog.LogData)
		lst.Add(*data)
	}

	err := container.Resolve[logData.Repository]().Save(lst)
	flog.ErrorIfExists(err)
}
