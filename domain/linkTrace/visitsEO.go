package linkTrace

import (
	"time"

	"github.com/farseer-go/fs/trace/eumTraceType"
)

type VisitsEO struct {
	AppName          string            // 应用名称
	VisitsNodePrefix string            // 前缀树访问节点
	VisitsNode       string            // 访问节点
	TraceType        eumTraceType.Enum // 入口类型
	MinMs            float64           // 最小访问速度
	MaxMs            float64           // 最大访问速度
	AvgMs            float64           // 平均访问速度
	Line95Ms         float64           // 95线访问速度
	Line99Ms         float64           // 99线访问速度
	ErrorCount       int               // 错误数量
	TotalCount       int               /// 总的数量
	QPS              float64           // 平均并发{
	CreateAt         time.Time         // 请求时间
}
