package localQueue

import (
	"fops/domain/logData"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
	"time"
)

func SaveFlogQueue(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	lst := mapper.ToList[flog.LogData](lstMessage)
	err := container.Resolve[logData.Repository]().Save(lst)
	flog.ErrorIfExists(err)
	time.Sleep(5 * time.Second)
}
