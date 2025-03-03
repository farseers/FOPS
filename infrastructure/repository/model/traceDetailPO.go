package model

import (
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/fs/trace/eumCallType"
)

type TraceDetailPO struct {
	TraceId        string                `gorm:"not null;default:'';comment:上下文ID"`
	AppId          string                `gorm:"not null;default:'';comment:应用ID"`
	AppName        string                `gorm:"not null;default:'';comment:应用名称"`
	AppIp          string                `gorm:"not null;default:'';comment:应用IP"`
	ParentAppName  string                `gorm:"not null;default:'';comment:上游应用"`
	DetailId       string                `gorm:"not null;default:0;comment:明细ID"`
	ParentDetailId string                `gorm:"not null;default:0;comment:父级明细ID"`
	Level          int                   `gorm:"not null;comment:当前层级（入口为0层）"`
	Comment        string                `gorm:"not null;default:'';comment:调用注释"`
	MethodName     string                `gorm:"not null;default:'';comment:调用方法"`
	CallType       eumCallType.Enum      `gorm:"not null;comment:调用类型"`
	Timeline       time.Duration         `gorm:"not null;default:0;comment:从入口开始统计（微秒）"`
	UnTraceTs      time.Duration         `gorm:"not null;default:0;comment:上一次结束到现在开始之间未Trace的时间（微秒）"`
	StartTs        int64                 `gorm:"not null;default:0;comment:调用开始时间戳（微秒）"`
	EndTs          int64                 `gorm:"not null;default:0;comment:调用停止时间戳（微秒）"`
	Exception      *trace.ExceptionStack `gorm:"type:String;json;not null;comment:异常信息"`
	CreateAt       dateTime.DateTime     `gorm:"type:DateTime64(3);not null;comment:请求时间"`

	TraceDetailHandPO     `gorm:"embedded;not null;comment:手动埋点" es_type:"object"`
	TraceDetailDatabasePO `gorm:"embedded;not null;comment:数据库埋点" es_type:"object"`
	TraceDetailEsPO       `gorm:"embedded;not null;comment:es埋点" es_type:"object"`
	TraceDetailEtcdPO     `gorm:"embedded;not null;comment:etcd埋点" es_type:"object"`
	TraceDetailEventPO    `gorm:"embedded;not null;comment:事件埋点" es_type:"object"`
	TraceDetailGrpcPO     `gorm:"embedded;not null;comment:grpc埋点" es_type:"object"`
	TraceDetailHttpPO     `gorm:"embedded;not null;comment:http埋点" es_type:"object"`
	TraceDetailMqPO       `gorm:"embedded;not null;comment:mq埋点" es_type:"object"`
	TraceDetailRedisPO    `gorm:"embedded;not null;comment:redis埋点" es_type:"object"`
	UseTs                 time.Duration `gorm:"not null;default:0;comment:总共使用时间微秒"`
	UseDesc               string        `gorm:"not null;default:'';comment:总共使用时间（描述）"`
}

// TraceDetailHandPO 手动埋点
type TraceDetailHandPO struct {
	HandName string `gorm:"not null;default:'';comment:名称"`
}

// TraceDetailDatabasePO 数据库埋点
type TraceDetailDatabasePO struct {
	DbName             string `gorm:"not null;default:'';comment:数据库名"`
	DbTableName        string `gorm:"not null;default:'';comment:表名"`
	DbSql              string `gorm:"not null;default:'';comment:SQL"`
	DbConnectionString string `gorm:"not null;default:'';comment:连接字符串"`
	DbRowsAffected     int64  `gorm:"not null;default:0;comment:影响行数"`
}

// TraceDetailEsPO es埋点
type TraceDetailEsPO struct {
	EsIndexName   string `gorm:"not null;default:'';comment:索引名称"`
	EsAliasesName string `gorm:"not null;default:'';comment:别名"`
}

// TraceDetailEtcdPO etcd埋点
type TraceDetailEtcdPO struct {
	EtcdKey     string `gorm:"not null;default:'';comment:etcd key"`
	EtcdLeaseID int64  `gorm:"not null;default:0;comment:LeaseID"`
}

// TraceDetailHandPO 事件埋点
type TraceDetailEventPO struct {
	EventName string `gorm:"not null;default:'';comment:名称"`
}

// TraceDetailHttpPO http埋点
type TraceDetailHttpPO struct {
	HttpMethod       string                                 `gorm:"not null;default:'';comment:post/get/put/delete"`
	HttpUrl          string                                 `gorm:"not null;default:'';comment:请求url"`
	HttpHeaders      collections.Dictionary[string, string] `gorm:"type:String;json;not null;comment:请求头部"`
	HttpRequestBody  string                                 `gorm:"not null;default:'';comment:入参"`
	HttpResponseBody string                                 `gorm:"not null;default:'';comment:出参"`
	HttpStatusCode   int                                    `gorm:"not null;default:0;comment:状态码"`
}

// TraceDetailGrpcPO grpc埋点
type TraceDetailGrpcPO struct {
	GrpcMethod       string                                 `gorm:"not null;default:'';comment:post/get/put/delete"`
	GrpcUrl          string                                 `gorm:"not null;default:'';comment:请求url"`
	GrpcHeaders      collections.Dictionary[string, string] `gorm:"type:String;json;not null;comment:请求头部"`
	GrpcRequestBody  string                                 `gorm:"not null;default:'';comment:入参"`
	GrpcResponseBody string                                 `gorm:"not null;default:'';comment:出参"`
	GrpcStatusCode   int                                    `gorm:"not null;default:0;comment:状态码"`
}

// TraceDetailMqPO mq埋点
type TraceDetailMqPO struct {
	MqServer     string `gorm:"not null;default:'';comment:MQ服务器地址"`
	MqExchange   string `gorm:"not null;default:'';comment:交换器名称"`
	MqRoutingKey string `gorm:"not null;default:'';comment:路由key"`
}

// TraceDetailRedisPO redis埋点
type TraceDetailRedisPO struct {
	RedisKey          string `gorm:"not null;default:'';comment:redis key"`
	RedisField        string `gorm:"not null;default:'';comment:hash field"`
	RedisRowsAffected int    `gorm:"not null;default:0;comment:影响行数"`
}
