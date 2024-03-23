package localQueue

import (
	"fops/domain/linkTrace"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	linkTraceCom "github.com/farseer-go/linkTrace"
	"time"
)

func SaveLinkTraceQueue(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	var lst collections.List[linkTraceCom.TraceContext]
	lstMessage.Select(&lst, func(item any) any {
		return item.(linkTraceCom.TraceContext)
	})

	err := container.Resolve[linkTrace.Repository]().Save(lst)
	flog.ErrorIfExists(err)
	time.Sleep(5 * time.Second)
}
