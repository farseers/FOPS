package linkTrace

import (
	"github.com/farseer-go/collections"
	linkTraceCom "github.com/farseer-go/linkTrace"
	"time"
)

// Repository 仓储接口
type Repository interface {
	// ToEntity 获取列表
	ToEntity(traceId string) collections.List[linkTraceCom.TraceContext]
	ToWebApiList(traceId, appName, appIp, requestIp, searchUrl string, statusCode int, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceContext]
	ToWebApiListByVisits(startAt, endAt time.Time) collections.List[linkTraceCom.TraceContext]
	ToTaskList(traceId, appName, appIp, taskName string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceContext]
	ToFScheduleList(traceId, appName, appIp, taskName string, taskGroupId, taskId, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceContext]
	ToConsumerList(traceId, appName, appIp, server, queueName, routingKey string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceContext]

	ToSlowDbList(traceId, appName, appIp, dbName, tableName string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailDatabase]
	ToSlowEsList(traceId, appName, appIp, indexName, aliasesName string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailEs]
	ToSlowEtcdList(traceId, appName, appIp, key string, leaseID, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailEtcd]
	ToSlowHandList(traceId, appName, appIp, name string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailHand]
	ToSlowHttpList(traceId, appName, appIp, method, url, requestBody, responseBody string, statusCode int, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailHttp]
	ToSlowMqList(traceId, appName, appIp, server, exchange, routingKey string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailMq]
	ToSlowRedisList(traceId, appName, appIp, key, field string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[linkTraceCom.TraceDetailRedis]
	Save(lstEO collections.List[linkTraceCom.TraceContext]) error

	SaveVisitsWebApi(lst collections.List[WebapiVisitsEO]) (int64, error)
	GetLastVisitsWebApiAt() (time.Time, error)
	ToWebApiVisitsList(appName, visitsNode string, startAt, endAt time.Time) collections.List[WebapiVisitsEO]
}
