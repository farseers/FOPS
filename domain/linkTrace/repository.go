package linkTrace

import (
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/fs/trace/eumTraceType"
	linkTraceCom "github.com/farseer-go/linkTrace"
)

// Repository 仓储接口
type Repository interface {
	// ToList 获取列表
	ToList(traceId string) collections.List[trace.TraceContext]
	Delete(traceType eumTraceType.Enum, startTime time.Time) error
	ToWebApiList(traceId, appName, appIp, requestIp, searchUrl string, statusCode int, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext]
	ToWebSocketList(traceId, appName, appIp, requestIp, searchUrl string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext]
	ToTraceListByVisits(startAt, endAt time.Time) collections.List[trace.TraceContext]
	ToTaskList(traceId, appName, appIp, taskName string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext]
	ToFScheduleList(traceId, appName, appIp, taskName string, taskGroupId, taskId, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext]
	ToConsumerList(traceId, appName, appIp, server, queueName, routingKey string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext]
	ToEventList(traceId, appName, appIp, queueName, routingKey string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext]
	ToQueueList(traceId, appName, appIp, queueName, routingKey string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext]

	DeleteSlow(dbName string, startTime time.Time) error

	ToSlowDbList(traceId, appName, appIp, dbName, tableName string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailDatabase]
	ToSlowEsList(traceId, appName, appIp, indexName, aliasesName string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailEs]
	ToSlowEtcdList(traceId, appName, appIp, key string, leaseID, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailEtcd]
	ToSlowHandList(traceId, appName, appIp, name string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[trace.TraceDetailHand]
	ToSlowHttpList(traceId, appName, appIp, method, url, body string, statusCode int, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailHttp]
	ToSlowMqList(traceId, appName, appIp, server, exchange, routingKey string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailMq]
	ToSlowRedisList(traceId, appName, appIp, key, field string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailRedis]
	Save(lstEO collections.List[trace.TraceContext]) error

	SaveVisits(lst collections.List[VisitsEO]) (int64, error)
	GetLastVisitsAt() (time.Time, error)
	ToVisitsList(appName, visitsNode string, startAt, endAt time.Time) collections.List[VisitsEO]
}
