package localQueue

import (
	"fops/domain/logData"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
)

func SaveFlogQueue(subscribeName string, lstMessage collections.ListAny, remainingCount int) {
	var lst collections.List[flog.LogData]
	lstMessage.Select(&lst, func(item any) any {
		return item.(flog.LogData)
	})

	err := container.Resolve[logData.Repository]().Save(lst)
	flog.ErrorIfExists(err)
}
