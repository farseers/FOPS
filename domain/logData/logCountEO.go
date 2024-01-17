package logData

import "github.com/farseer-go/fs/core/eumLogLevel"

type LogCountEO struct {
	LogCount int              // 日志数量
	LogType  eumLogLevel.Enum // 日志等级
}
