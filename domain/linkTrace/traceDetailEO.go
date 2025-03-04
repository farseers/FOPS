package linkTrace

import (
	"time"

	"github.com/farseer-go/fs/trace"
)

// TraceDetail 埋点明细（基类）
type TraceDetailEO struct {
	trace.TraceDetail
	TraceId       string        // 上下文ID
	AppId         string        // 应用ID
	AppName       string        // 应用名称
	AppIp         string        // 应用IP
	ParentAppName string        // 上游应用
	UseTs         time.Duration // 总共使用时间微秒
	UseDesc       string        // 总共使用时间（描述）
}
