package logData

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/flog"
	"time"
)

// Repository 仓储接口
type Repository interface {
	Save(lstEO collections.List[flog.LogData]) error
	ToInfo(id string) flog.LogData
	ToList(traceId string, appName, appIp, logContent string, minute int, logLevel eumLogLevel.Enum, pageSize, pageIndex int) collections.PageList[flog.LogData]
	StatCount() collections.List[LogCountEO]
	DeleteLog(startTime time.Time) error
}
