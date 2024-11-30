package model

import (
	"time"

	"github.com/farseer-go/fs/trace/eumTraceType"
)

type VisitsPO struct {
	CreateAt         time.Time         `gorm:"type:DateTime64(3);not null;comment:请求时间"`
	AppName          string            `gorm:"not null;comment:应用名称"`
	VisitsNodePrefix string            `gorm:"not null;comment:前缀树访问节点"`
	VisitsNode       string            `gorm:"not null;comment:访问节点"`
	TraceType        eumTraceType.Enum `gorm:"not null;comment:入口类型"`
	MinMs            float64           `gorm:"type:Float64(12,3);not null;default:0;comment:最小访问速度"`
	MaxMs            float64           `gorm:"type:Float64(12,3);not null;default:0;comment:最大访问速度"`
	AvgMs            float64           `gorm:"type:Float64(12,3);not null;comment:平均访问速度"`
	Line95Ms         float64           `gorm:"type:Float64(12,3);not null;default:0;comment:95线访问速度"`
	Line99Ms         float64           `gorm:"type:Float64(12,3);not null;default:0;comment:99线访问速度"`
	ErrorCount       int               `gorm:"type:Int64;not null;default:0;comment:错误数量"`
	TotalCount       int               `gorm:"type:Int64;not null;default:0;comment:总的数量"`
	QPS              float64           `gorm:"type:Float64(12,3);not null;comment:平均并发"`
}
