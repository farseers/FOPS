package localQueue

import (
	"fops/domain/linkTrace"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
)

func SaveLinkTraceQueue(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	lst := collections.NewList[trace.TraceContext]()
	for _, item := range lstMessage.ToArray() {
		data := item.(*trace.TraceContext)
		lst.Add(*data)
	}

	err := container.Resolve[linkTrace.Repository]().Save(lst)
	flog.ErrorIfExists(err)
}
