package localQueue

import (
	"fmt"
	"fops/domain/linkTrace"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	linkTraceCom "github.com/farseer-go/linkTrace"
)

func SaveLinkTraceQueue(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	var lst collections.List[linkTraceCom.TraceContext]
	lstMessage.Select(&lst, func(item any) any {
		return item.(linkTraceCom.TraceContext)
	})

	fmt.Printf("SaveLinkTraceQueue：准备保存%d条\n", lstMessage.Count())
	err := container.Resolve[linkTrace.Repository]().Save(lst)
	fmt.Printf("SaveLinkTraceQueue：保存结束\n")
	flog.ErrorIfExists(err)
}
