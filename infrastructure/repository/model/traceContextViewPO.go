package model

import (
	_ "embed"
	"time"

	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/trace/eumTraceType"
)

type TraceContextViewPO struct {
	TraceId           string            `gorm:"not null;default:'';comment:上下文ID"`
	AppId             string            `gorm:"not null;default:'';comment:应用ID"`
	AppName           string            `gorm:"not null;default:'';comment:应用名称"`
	AppIp             string            `gorm:"not null;default:'';comment:应用IP"`
	ParentAppName     string            `gorm:"not null;default:'';comment:上游应用"`
	TraceLevel        int               `gorm:"not null;default:0;comment:逐层递增（显示上下游顺序）"`
	TraceCount        int               `gorm:"not null;default:0;comment:追踪明细数量"`
	StartTs           int64             `gorm:"not null;default:0;comment:调用开始时间戳（微秒）"`
	EndTs             int64             `gorm:"not null;default:0;comment:调用结束时间戳（微秒）"`
	UseTs             time.Duration     `gorm:"not null;default:0;comment:总共使用时间（微秒）"`
	UseDesc           string            `gorm:"not null;default:'';comment:总共使用时间（描述）"`
	TraceType         eumTraceType.Enum `gorm:"not null;comment:入口类型"`
	Exception         *ExceptionStackPO `gorm:"json;not null;comment:异常信息"`
	WebContextPO      `gorm:"embedded;not null;comment:Web请求上下文" es_type:"object"`
	ConsumerContextPO `gorm:"embedded;not null;comment:消费上下文" es_type:"object"`
	TaskContextPO     `gorm:"embedded;not null;comment:任务上下文" es_type:"object"`
	WatchKeyContextPO `gorm:"embedded;not null;comment:Etcd上下文" es_type:"object"`
	CreateAt          dateTime.DateTime `gorm:"type:DateTime64(3);not null;comment:请求时间"`
}

//go:embed sql/v_linkTrace.sql
var vLinkTraceSql string

// CreateTable 创建表
func (*TraceContextViewPO) CreateTable() string {
	return vLinkTraceSql
}
