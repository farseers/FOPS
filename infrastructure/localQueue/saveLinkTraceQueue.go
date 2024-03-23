package localQueue

import (
	"fops/domain/linkTrace"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	linkTraceCom "github.com/farseer-go/linkTrace"
	"github.com/farseer-go/mapper"
	"time"
)

func SaveLinkTraceQueue(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	lst := mapper.ToList[linkTraceCom.TraceContext](lstMessage)
	err := container.Resolve[linkTrace.Repository]().Save(lst)
	flog.ErrorIfExists(err)
	time.Sleep(5 * time.Second)
}
