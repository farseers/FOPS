package localQueue

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/trace"
)

func SaveLinkTraceQueue(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	//lst := collections.NewList[linkTraceCom.TraceContext]()
	//lstMessage.Foreach(func(item *any) {
	//	lst.Add((*item).(linkTraceCom.TraceContext))
	//})
	//
	//err := container.Resolve[linkTrace.Repository]().Save(lst)
	//flog.ErrorIfExists(err)
}
