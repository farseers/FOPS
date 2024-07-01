package localQueue

import (
	"fmt"
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

	var lst collections.List[flog.LogData]
	lstMessage.Select(&lst, func(item any) any {
		return item.(flog.LogData)
	})
	fmt.Printf("SaveFlogQueue：准备保存%d条\n", lstMessage.Count())
	err := container.Resolve[logData.Repository]().Save(lst)
	fmt.Printf("SaveFlogQueue：保存结束\n")
	flog.ErrorIfExists(err)
}
