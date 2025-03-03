package repository

import (
	"bytes"
	"fmt"
	"fops/domain/linkTrace"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/fs/trace/eumCallType"
	"github.com/farseer-go/fs/trace/eumTraceType"
	"github.com/farseer-go/mapper"
)

type linkTraceRepository struct{}

func (receiver *linkTraceRepository) ToList(traceId string) collections.List[trace.TraceContext] {
	lst := collections.NewList[trace.TraceContext]()
	if linkTrace.Config.Driver == "clickhouse" {
		lstPO := context.CHContext.TraceContext.Where("trace_id", traceId).Asc("start_ts").ToList()
		lst = mapper.ToList[trace.TraceContext](lstPO)
	}
	return lst
}

func (receiver *linkTraceRepository) Delete(traceType eumTraceType.Enum, startTime time.Time) error {
	if linkTrace.Config.Driver == "clickhouse" {
		_, err := context.CHContext.TraceContext.Where("create_at <= ?", startTime).Delete() // trace_type = ? and
		return err
	}
	return nil
}

func (receiver *linkTraceRepository) ToWebApiList(traceId, appName, appIp, requestIp, searchUrl string, statusCode int, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceContext.Select("trace_id,app_id,app_name,app_ip,parent_app_name,trace_type,start_ts,end_ts,use_ts,use_desc,trace_count,create_at,exception,web_domain,web_path,web_method,web_content_type,web_status_code,web_request_ip").
			Where("trace_type = ? and parent_app_name = ''", eumTraceType.WebApi).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Millisecond)).
			WhereIf(requestIp != "", "web_request_ip = ?", requestIp).
			WhereIf(searchUrl != "", "web_path like ?", "%"+searchUrl+"%").
			WhereIf(statusCode > 0, "web_status_code = ?", statusCode).
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		return mapper.ToPageList[trace.TraceContext](lstPO)
	}
	return collections.NewPageList[trace.TraceContext](collections.NewList[trace.TraceContext](), 0)
}

func (receiver *linkTraceRepository) ToWebSocketList(traceId, appName, appIp, requestIp, searchUrl string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceContext.Select("trace_id,app_id,app_name,app_ip,parent_app_name,trace_type,start_ts,end_ts,use_ts,use_desc,trace_count,create_at,exception,web_domain,web_path,web_method,web_request_ip").
			Where("trace_type = ? and parent_app_name = ''", eumTraceType.WebSocket).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Millisecond)).
			WhereIf(requestIp != "", "web_request_ip = ?", requestIp).
			WhereIf(searchUrl != "", "web_path like ?", "%"+searchUrl+"%").
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		return mapper.ToPageList[trace.TraceContext](lstPO)
	}
	return collections.NewPageList[trace.TraceContext](collections.NewList[trace.TraceContext](), 0)
}
func (receiver *linkTraceRepository) ToTaskList(traceId, appName, appIp, taskName string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceContext.Select("trace_id,app_id,app_name,app_ip,parent_app_name,trace_type,start_ts,end_ts,use_ts,use_desc,trace_count,create_at,exception,task_name").
			Where("trace_type = ?", eumTraceType.Task).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Millisecond)).
			WhereIf(taskName != "", "task_name like ?", "%"+taskName+"%").
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		return mapper.ToPageList[trace.TraceContext](lstPO)
	}
	return collections.NewPageList[trace.TraceContext](collections.NewList[trace.TraceContext](), 0)
}
func (receiver *linkTraceRepository) ToFScheduleList(traceId, appName, appIp, taskName string, taskGroupId, taskId, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceContext.Select("trace_id,app_id,app_name,app_ip,parent_app_name,trace_type,start_ts,end_ts,use_ts,use_desc,trace_count,create_at,exception,task_name,task_group_name,task_id,task_data").
			Where("trace_type = ? and parent_app_name = ''", eumTraceType.FSchedule).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Millisecond)).
			WhereIf(taskName != "", "task_name like ?", "%"+taskName+"%").
			WhereIf(taskGroupId > 0, "task_group_id = ?", taskGroupId).
			WhereIf(taskId > 0, "task_id = ?", taskId).
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		return mapper.ToPageList[trace.TraceContext](lstPO)
	}
	return collections.NewPageList[trace.TraceContext](collections.NewList[trace.TraceContext](), 0)
}
func (receiver *linkTraceRepository) ToConsumerList(traceId, appName, appIp, server, queueName, routingKey string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceContext.Select("trace_id,app_id,app_name,app_ip,parent_app_name,trace_type,start_ts,end_ts,use_ts,use_desc,trace_count,create_at,exception,consumer_server,consumer_queue_name,consumer_routing_key").
			Where("trace_type = ?", eumTraceType.MqConsumer).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Millisecond)).
			WhereIf(server != "", "consumer_server like ?", "%"+server+"%").
			WhereIf(queueName != "", "consumer_queue_name like ?", "%"+queueName+"%").
			WhereIf(routingKey != "", "consumer_routing_key like ?", "%"+routingKey+"%").
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		return mapper.ToPageList[trace.TraceContext](lstPO)
	}
	return collections.NewPageList[trace.TraceContext](collections.NewList[trace.TraceContext](), 0)
}
func (receiver *linkTraceRepository) ToEventList(traceId, appName, appIp, queueName, routingKey string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceContext.Select("trace_id,app_id,app_name,app_ip,parent_app_name,trace_type,start_ts,end_ts,use_ts,use_desc,trace_count,create_at,exception,consumer_server,consumer_queue_name,consumer_routing_key").
			Where("trace_type = ?", eumTraceType.EventConsumer).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Millisecond)).
			WhereIf(queueName != "", "consumer_queue_name like ?", "%"+queueName+"%").
			WhereIf(routingKey != "", "consumer_routing_key like ?", "%"+routingKey+"%").
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		return mapper.ToPageList[trace.TraceContext](lstPO)
	}
	return collections.NewPageList[trace.TraceContext](collections.NewList[trace.TraceContext](), 0)
}
func (receiver *linkTraceRepository) ToQueueList(traceId, appName, appIp, queueName, routingKey string, searchUseTs int64, onlyViewException bool, startMin int, pageSize, pageIndex int) collections.PageList[trace.TraceContext] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceContext.Select("trace_id,app_id,app_name,app_ip,parent_app_name,trace_type,start_ts,end_ts,use_ts,use_desc,trace_count,create_at,exception,consumer_server,consumer_queue_name,consumer_routing_key").
			Where("trace_type = ?", eumTraceType.QueueConsumer).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Millisecond)).
			WhereIf(queueName != "", "consumer_queue_name like ?", "%"+queueName+"%").
			WhereIf(routingKey != "", "consumer_routing_key like ?", "%"+routingKey+"%").
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		return mapper.ToPageList[trace.TraceContext](lstPO)
	}
	return collections.NewPageList[trace.TraceContext](collections.NewList[trace.TraceContext](), 0)
}

func (receiver *linkTraceRepository) ToTraceListByVisits(startAt, endAt time.Time) collections.List[trace.TraceContext] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceContext.Asc("use_ts").
			Where("start_ts >= ? and start_ts < ?", startAt.UnixMicro(), endAt.UnixMicro()) // parent_app_name = '' and
		lstPO := ts.Asc("use_ts").ToList()
		return mapper.ToList[trace.TraceContext](lstPO)
		//var arr []trace.TraceContext
		//ts.Fill(&arr)
		//return collections.NewList(arr...)
	}
	return collections.NewList[trace.TraceContext]()
}

func (receiver *linkTraceRepository) DeleteSlow(dbName string, startTime time.Time) error {
	if linkTrace.Config.Driver == "clickhouse" {
		context.CHContext.TraceDetail.Where("create_at <= ?", startTime).Delete()
	}
	return nil
}

func (receiver *linkTraceRepository) ToSlowDbList(traceId, appName, appIp, dbName, tableName string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[trace.TraceDetail] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceDetail.
			Where("call_type = ?", eumCallType.Database).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Microsecond)).
			WhereIf(dbName != "", "db_name like ?", "%"+dbName+"%").
			WhereIf(tableName != "", "table_name like ?", "%"+tableName+"%").
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		lst := mapper.ToPageList[trace.TraceDetail](lstPO)
		return lst
	}
	return collections.NewPageList[trace.TraceDetail](collections.NewList[trace.TraceDetail](), 0)
}
func (receiver *linkTraceRepository) ToSlowEsList(traceId, appName, appIp, indexName, aliasesName string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[trace.TraceDetail] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceDetail.
			Where("call_type = ?", eumCallType.Elasticsearch).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Microsecond)).
			WhereIf(indexName != "", "index_name like ?", "%"+indexName+"%").
			WhereIf(aliasesName != "", "aliases_name like ?", "%"+aliasesName+"%").
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		lst := mapper.ToPageList[trace.TraceDetail](lstPO)
		return lst
	}
	return collections.NewPageList[trace.TraceDetail](collections.NewList[trace.TraceDetail](), 0)
}
func (receiver *linkTraceRepository) ToSlowEtcdList(traceId, appName, appIp, key string, leaseID, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[trace.TraceDetail] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceDetail.
			Where("call_type = ?", eumCallType.Etcd).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Microsecond)).
			WhereIf(key != "", "key like ?", "%"+key+"%").
			WhereIf(leaseID > 0, "leaseID = ?", leaseID).
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		lst := mapper.ToPageList[trace.TraceDetail](lstPO)
		return lst
	}
	return collections.NewPageList[trace.TraceDetail](collections.NewList[trace.TraceDetail](), 0)
}
func (receiver *linkTraceRepository) ToSlowHandList(traceId, appName, appIp, name string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[trace.TraceDetail] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceDetail.
			Where("call_type = ?", eumCallType.Hand).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Microsecond)).
			WhereIf(name != "", "name like ?", "%"+name+"%").
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		lst := mapper.ToPageList[trace.TraceDetail](lstPO)
		return lst
	}
	return collections.NewPageList[trace.TraceDetail](collections.NewList[trace.TraceDetail](), 0)
}
func (receiver *linkTraceRepository) ToSlowHttpList(traceId, appName, appIp, method, url, body string, statusCode int, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[trace.TraceDetail] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceDetail.
			Where("call_type = ?", eumCallType.Http).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Microsecond)).
			WhereIf(method != "", "method like ?", "%"+method+"%").
			WhereIf(url != "", "url like ?", "%"+url+"%").
			WhereIf(body != "", "(request_body like ? or response_body like ?)", "%"+body+"%", "%"+body+"%").
			WhereIf(statusCode > 0, "status_code >= ?", statusCode).
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		lst := mapper.ToPageList[trace.TraceDetail](lstPO)
		return lst
	}
	return collections.NewPageList[trace.TraceDetail](collections.NewList[trace.TraceDetail](), 0)
}
func (receiver *linkTraceRepository) ToSlowMqList(traceId, appName, appIp, server, exchange, routingKey string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[trace.TraceDetail] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceDetail.
			Where("call_type = ?", eumCallType.Mq).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Microsecond)).
			WhereIf(server != "", "server like ?", "%"+server+"%").
			WhereIf(exchange != "", "exchange like ?", "%"+exchange+"%").
			WhereIf(routingKey != "", "url like ?", "%"+routingKey+"%").
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		lst := mapper.ToPageList[trace.TraceDetail](lstPO)
		return lst
	}
	return collections.NewPageList[trace.TraceDetail](collections.NewList[trace.TraceDetail](), 0)
}
func (receiver *linkTraceRepository) ToSlowRedisList(traceId, appName, appIp, key, field string, searchUseTs int64, onlyViewException bool, startMin, pageSize, pageIndex int) collections.PageList[trace.TraceDetail] {
	if linkTrace.Config.Driver == "clickhouse" {
		ts := context.CHContext.TraceDetail.
			Where("call_type = ?", eumCallType.Redis).
			WhereIf(traceId != "", "trace_id = ?", traceId).
			WhereIf(appName != "", "LOWER(app_name) = ?", appName).
			WhereIf(appIp != "", "app_ip = ?", appIp).
			WhereIf(searchUseTs > 0, "use_ts >= ?", searchUseTs*int64(time.Microsecond)).
			WhereIf(key != "", "key like ?", "%"+key+"%").
			WhereIf(field != "", "field like ?", "%"+field+"%").
			WhereIf(onlyViewException, "exception <> ''").
			WhereIf(startMin > 0, "start_ts >= ?", dateTime.Now().AddMinutes(-startMin).UnixMicro())

		lstPO := ts.DescIfElse(startMin > 0, "use_ts", "start_ts").ToPageList(pageSize, pageIndex)
		lst := mapper.ToPageList[trace.TraceDetail](lstPO)
		return lst
	}
	return collections.NewPageList[trace.TraceDetail](collections.NewList[trace.TraceDetail](), 0)
}

func (receiver *linkTraceRepository) Save(lstEO collections.List[trace.TraceContext]) error {
	lst := collections.NewList[model.TraceContextPO]()
	lstDetail := collections.NewList[model.TraceDetailPO]()
	lstEO.Foreach(func(item *trace.TraceContext) {
		po := model.TraceContextPO{
			TraceId:       item.TraceId,
			AppId:         item.AppId,
			AppName:       item.AppName,
			AppIp:         item.AppIp,
			ParentAppName: item.ParentAppName,
			TraceLevel:    item.TraceLevel,
			TraceCount:    item.TraceCount,
			StartTs:       item.StartTs,
			EndTs:         item.EndTs,
			UseTs:         item.UseTs,
			UseDesc:       item.UseDesc,
			TraceType:     item.TraceType,
			List:          item.List,
			WebContextPO: model.WebContextPO{
				WebDomain:       item.WebContext.WebDomain,
				WebPath:         item.WebContext.WebPath,
				WebMethod:       item.WebContext.WebMethod,
				WebContentType:  item.WebContext.WebContentType,
				WebStatusCode:   item.WebContext.WebStatusCode,
				WebHeaders:      item.WebContext.WebHeaders,
				WebRequestBody:  item.WebContext.WebRequestBody,
				WebResponseBody: item.WebContext.WebResponseBody,
				WebRequestIp:    item.WebContext.WebRequestIp,
			},
			ConsumerContextPO: model.ConsumerContextPO{
				ConsumerServer:     item.ConsumerContext.ConsumerServer,
				ConsumerQueueName:  item.ConsumerContext.ConsumerQueueName,
				ConsumerRoutingKey: item.ConsumerContext.ConsumerRoutingKey,
			},
			TaskContextPO: model.TaskContextPO{
				TaskName:      item.TaskContext.TaskName,
				TaskGroupName: item.TaskContext.TaskGroupName,
				TaskId:        item.TaskContext.TaskId,
				TaskData:      item.TaskContext.TaskData,
			},
			WatchKeyContextPO: model.WatchKeyContextPO{
				WatchKey: item.WatchKeyContext.WatchKey,
			},
			CreateAt: item.CreateAt,
		}
		if item.Exception != nil {
			po.Exception = &model.ExceptionStackPO{
				ExceptionCallFile:     item.Exception.ExceptionCallFile,
				ExceptionCallLine:     item.Exception.ExceptionCallLine,
				ExceptionCallFuncName: item.Exception.ExceptionCallFuncName,
				ExceptionIsException:  item.Exception.ExceptionIsException,
				ExceptionMessage:      item.Exception.ExceptionMessage,
			}
		}
		lst.Add(po)

		// 转换明细
		item.List.Foreach(func(detail **trace.TraceDetail) {
			traceDetail := *detail
			lstDetail.Add(model.TraceDetailPO{
				TraceId:        item.TraceId,
				AppId:          item.AppId,
				AppName:        item.AppName,
				AppIp:          item.AppIp,
				ParentAppName:  item.ParentAppName,
				DetailId:       traceDetail.DetailId,
				ParentDetailId: traceDetail.ParentDetailId,
				Level:          traceDetail.Level,
				Comment:        traceDetail.Comment,
				MethodName:     traceDetail.MethodName,
				CallType:       traceDetail.CallType,
				Timeline:       traceDetail.Timeline,
				UnTraceTs:      traceDetail.UnTraceTs,
				StartTs:        traceDetail.StartTs,
				EndTs:          traceDetail.EndTs,
				Exception:      traceDetail.Exception,
				CreateAt:       traceDetail.CreateAt,
				TraceDetailHandPO: model.TraceDetailHandPO{
					HandName: traceDetail.HandName,
				},
				TraceDetailDatabasePO: model.TraceDetailDatabasePO{
					DbName:             traceDetail.DbName,
					DbTableName:        traceDetail.DbTableName,
					DbSql:              traceDetail.DbSql,
					DbConnectionString: traceDetail.DbConnectionString,
					DbRowsAffected:     traceDetail.DbRowsAffected,
				},
				TraceDetailEsPO: model.TraceDetailEsPO{
					EsIndexName:   traceDetail.EsIndexName,
					EsAliasesName: traceDetail.EsAliasesName,
				},
				TraceDetailEtcdPO: model.TraceDetailEtcdPO{
					EtcdKey:     traceDetail.EtcdKey,
					EtcdLeaseID: traceDetail.EtcdLeaseID,
				},
				TraceDetailEventPO: model.TraceDetailEventPO{
					EventName: traceDetail.EventName,
				},
				TraceDetailGrpcPO: model.TraceDetailGrpcPO{
					GrpcMethod:       traceDetail.GrpcMethod,
					GrpcUrl:          traceDetail.GrpcUrl,
					GrpcHeaders:      traceDetail.GrpcHeaders,
					GrpcRequestBody:  traceDetail.GrpcRequestBody,
					GrpcResponseBody: traceDetail.GrpcResponseBody,
					GrpcStatusCode:   traceDetail.GrpcStatusCode,
				},
				TraceDetailHttpPO: model.TraceDetailHttpPO{
					HttpMethod:       traceDetail.HttpMethod,
					HttpUrl:          traceDetail.HttpUrl,
					HttpHeaders:      traceDetail.HttpHeaders,
					HttpRequestBody:  traceDetail.HttpRequestBody,
					HttpResponseBody: traceDetail.HttpResponseBody,
					HttpStatusCode:   traceDetail.HttpStatusCode,
				},
				TraceDetailMqPO: model.TraceDetailMqPO{
					MqServer:     traceDetail.MqServer,
					MqExchange:   traceDetail.MqExchange,
					MqRoutingKey: traceDetail.MqRoutingKey,
				},
				TraceDetailRedisPO: model.TraceDetailRedisPO{
					RedisKey:          traceDetail.RedisKey,
					RedisField:        traceDetail.RedisField,
					RedisRowsAffected: traceDetail.RedisRowsAffected,
				},
				UseTs:   time.Duration(traceDetail.EndTs-traceDetail.StartTs) * time.Microsecond,
				UseDesc: (time.Duration(traceDetail.EndTs-traceDetail.StartTs) * time.Microsecond).String(),
			})
		})
	})

	if linkTrace.Config.Driver == "clickhouse" {
		// 写入上下文
		if _, err := context.CHContext.TraceContext.InsertList(lst, 2000); err != nil {
			_ = flog.Errorf("TraceContext写入ch失败,%s", err.Error())
		}
	} else {
		return fmt.Errorf("不支持的链路追踪驱动：%s", linkTrace.Config.Driver)
	}

	// 写入明细
	return receiver.saveDetail(lstDetail)
}

// 写入明细
func (receiver *linkTraceRepository) saveDetail(lstDetail collections.List[model.TraceDetailPO]) error {
	if linkTrace.Config.Driver == "clickhouse" {
		if _, err := context.CHContext.TraceDetail.InsertList(lstDetail, 2000); err != nil {
			return err
		}
	}
	return nil
}

func (receiver *linkTraceRepository) SaveVisits(lst collections.List[linkTrace.VisitsEO]) (int64, error) {
	lstPO := mapper.ToList[model.VisitsPO](lst)
	if linkTrace.Config.Driver == "clickhouse" {
		return context.CHContext.Visits.InsertList(lstPO, 10000)
	}

	return 0, fmt.Errorf("不支持的链路追踪驱动：%s", linkTrace.Config.Driver)
}

func (receiver *linkTraceRepository) GetLastVisitsAt() (time.Time, error) {
	if linkTrace.Config.Driver == "clickhouse" {
		return context.CHContext.Visits.Desc("create_at").GetTime("create_at"), nil
	}

	return time.Time{}, fmt.Errorf("不支持的链路追踪驱动：%s", linkTrace.Config.Driver)
}

func (receiver *linkTraceRepository) ToVisitsList(appName, visitsNode string, startAt, endAt time.Time) collections.List[linkTrace.VisitsEO] {
	if linkTrace.Config.Driver == "clickhouse" {
		sql := bytes.Buffer{}
		sql.WriteString("select visits_node,min(min_ms) as min_ms,max(max_ms) as max_ms,avg(avg_ms) as avg_ms,avg(line95_ms) as line95_ms,avg(line99_ms) as line99_ms,sum(error_count) as error_count,sum(total_count) as total_count ,max(qps) as qps from visits ")
		sql.WriteString(fmt.Sprintf("where visits_node_prefix = '%s' and create_at >= '%s' and create_at < '%s' ", visitsNode, startAt.Format(time.DateTime), endAt.Format(time.DateTime)))
		if appName != "" {
			sql.WriteString(fmt.Sprintf("and app_name = '%s' ", appName))
		}
		sql.WriteString("group by visits_node ")
		sql.WriteString("order by visits_node asc")

		lstPO := context.CHContext.Visits.ExecuteSqlToList(sql.String())
		return mapper.ToList[linkTrace.VisitsEO](lstPO)
	}
	return collections.NewList[linkTrace.VisitsEO]()
}
