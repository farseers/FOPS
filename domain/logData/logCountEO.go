package logData

import "github.com/farseer-go/fs/core/eumLogLevel"

type LogCountEO struct {
	AppName  string           // 应用名称
	LogLevel eumLogLevel.Enum // 日志等级
	LogCount int              // 日志数量
}
