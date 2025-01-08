package model

import (
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/fs/trace/eumCallType"
)

type BaseTraceDetailPO struct {
	TraceId        string                `gorm:"not null;default:'';comment:上下文ID"`
	AppId          string                `gorm:"not null;default:'';comment:应用ID"`
	AppName        string                `gorm:"not null;default:'';comment:应用名称"`
	AppIp          string                `gorm:"not null;default:'';comment:应用IP"`
	ParentAppName  string                `gorm:"not null;default:'';comment:上游应用"`
	DetailId       string                `gorm:"not null;default:0;comment:明细ID"`
	ParentDetailId string                `gorm:"not null;default:0;comment:父级明细ID"`
	Level          int                   `gorm:"not null;comment:当前层级（入口为0层）"`
	CallType       eumCallType.Enum      `gorm:"not null;comment:调用类型"`
	Timeline       time.Duration         `gorm:"not null;default:0;comment:从入口开始统计（微秒）"`
	UnTraceTs      time.Duration         `gorm:"not null;default:0;comment:上一次结束到现在开始之间未Trace的时间（微秒）"`
	StartTs        int64                 `gorm:"not null;default:0;comment:调用开始时间戳（微秒）"`
	EndTs          int64                 `gorm:"not null;default:0;comment:调用停止时间戳（微秒）"`
	UseTs          time.Duration         `gorm:"not null;default:0;comment:总共使用时间微秒"`
	UseDesc        string                `gorm:"not null;default:'';comment:总共使用时间（描述）"`
	Exception      *trace.ExceptionStack `gorm:"json;not null;comment:异常信息"`
	MethodName     string                `gorm:"not null;default:'';comment:调用方法"`
	Comment        string                `gorm:"not null;default:'';comment:调用注释"`
	CreateAt       dateTime.DateTime     `gorm:"type:DateTime64(3);not null;comment:请求时间"`
}

type TraceDetailDatabasePO struct {
	BaseTraceDetailPO `gorm:"embedded"`
	DbName            string `gorm:"not null;default:'';comment:数据库名"`
	TableName         string `gorm:"not null;default:'';comment:表名"`
	Sql               string `gorm:"not null;default:'';comment:SQL"`
	ConnectionString  string `gorm:"not null;default:'';comment:连接字符串"`
	RowsAffected      int64  `gorm:"not null;default:0;comment:影响行数"`
}

type TraceDetailEsPO struct {
	BaseTraceDetailPO `gorm:"embedded"`
	IndexName         string `gorm:"not null;default:'';comment:索引名称"`
	AliasesName       string `gorm:"not null;default:'';comment:别名"`
}
type TraceDetailEtcdPO struct {
	BaseTraceDetailPO `gorm:"embedded"`
	Key               string `gorm:"not null;default:'';comment:etcd key"`
	LeaseID           int64  `gorm:"not null;default:0;comment:LeaseID"`
}

// TraceDetailHandPO 手动埋点
type TraceDetailHandPO struct {
	BaseTraceDetailPO `gorm:"embedded"`
	Name              string `gorm:"not null;default:'';comment:名称"`
}
type TraceDetailHttpPO struct {
	BaseTraceDetailPO `gorm:"embedded"`
	Method            string                                 `gorm:"not null;default:'';comment:post/get/put/delete"`
	Url               string                                 `gorm:"not null;default:'';comment:请求url"`
	Headers           collections.Dictionary[string, string] `gorm:"type:String;json;not null;comment:请求头部"`
	RequestBody       string                                 `gorm:"not null;default:'';comment:入参"`
	ResponseBody      string                                 `gorm:"not null;default:'';comment:出参"`
	StatusCode        int                                    `gorm:"not null;default:0;comment:状态码"`
}
type TraceDetailGrpcPO struct {
	BaseTraceDetailPO `gorm:"embedded"`
	Method            string                                 `gorm:"not null;default:'';comment:post/get/put/delete"`
	Url               string                                 `gorm:"not null;default:'';comment:请求url"`
	Headers           collections.Dictionary[string, string] `gorm:"type:String;json;not null;comment:请求头部"`
	RequestBody       string                                 `gorm:"not null;default:'';comment:入参"`
	ResponseBody      string                                 `gorm:"not null;default:'';comment:出参"`
	StatusCode        int                                    `gorm:"not null;default:0;comment:状态码"`
}
type TraceDetailMqPO struct {
	BaseTraceDetailPO `gorm:"embedded"`
	Server            string `gorm:"not null;default:'';comment:MQ服务器地址"`
	Exchange          string `gorm:"not null;default:'';comment:交换器名称"`
	RoutingKey        string `gorm:"not null;default:'';comment:路由key"`
}
type TraceDetailRedisPO struct {
	BaseTraceDetailPO `gorm:"embedded"`
	Key               string `gorm:"not null;default:'';comment:redis key"`
	Field             string `gorm:"not null;default:'';comment:hash field"`
	RowsAffected      int    `gorm:"not null;default:0;comment:影响行数"`
}
