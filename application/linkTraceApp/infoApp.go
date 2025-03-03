// @area /linkTrace/
package linkTraceApp

import (
	"fmt"
	"fops/application/linkTraceApp/response"
	"fops/domain/linkTrace"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/fs/trace/eumCallType"
	"github.com/farseer-go/fs/trace/eumTraceType"
)

// Info 链路追踪日志详情
// @get info/{traceId}
func Info(traceId string, linkTraceRepository linkTrace.Repository) response.LinkTraceResponse {
	l := linkTraceWarp{
		lst:       collections.NewList[response.LinkTraceVO](),
		rgbaIndex: -1,
		lstPO:     linkTraceRepository.ToList(traceId),
	}
	if l.lstPO.Count() == 0 {
		return response.LinkTraceResponse{}
	}

	// 找到最早的开始时间 和 最晚的结束时间，得到总耗时
	l.startTs = l.lstPO.Min(func(item trace.TraceContext) any {
		return item.StartTs
	}).(int64)
	l.endTs = l.lstPO.Max(func(item trace.TraceContext) any {
		return item.EndTs
	}).(int64)
	l.TotalUse = float64(l.endTs - l.startTs)

	// 当A服务调用B服务时，前后均有可能包含数据库之类的操作。因此需要将lstPO重新组织。按实际的调用顺序重新排序
	// 前端就可以简单的遍历lst显示到页面即可
	entryPO := l.lstPO.Where(func(item trace.TraceContext) bool {
		return item.ParentAppName == "" // 不同服务的机器时间会有差异，不能直接通过start_ts来排序
	}).First()
	// 没有取到，则根据start_ts排序，取第一个
	if entryPO.TraceId == "" {
		entryPO = l.lstPO.OrderBy(func(item trace.TraceContext) any {
			return item.StartTs
		}).First()
	}

	l.addEntry(entryPO)

	// 补充剩余的
	for l.lstPO.Count() > 0 {
		entryPO2 := l.lstPO.OrderBy(func(item trace.TraceContext) any {
			return item.StartTs
		}).First()
		l.addEntry(entryPO2)
	}

	// 调用完成
	l.lst.Add(response.LinkTraceVO{
		Rgba: response.RgbaList[0], AppId: entryPO.AppId, AppIp: entryPO.AppIp, AppName: entryPO.AppName, UseTs: 0,
		StartTs:   float64(entryPO.EndTs - entryPO.StartTs),
		StartRate: l.lst.Last().StartRate,
		UseRate:   0,
		Caption:   "调用完成",
		Desc:      fmt.Sprintf("耗时：%v ms", entryPO.UseTs.Milliseconds()),
	})

	rsp := response.LinkTraceResponse{
		Entry: entryPO,
		List:  l.lst,
	}
	rsp.Entry.UseDesc = rsp.Entry.UseTs.String()
	rsp.Entry.List.Clear()
	rsp.Entry.StartTs = l.startTs
	rsp.Entry.EndTs = l.endTs
	rsp.Entry.UseTs = time.Duration(l.TotalUse) * time.Microsecond
	rsp.Entry.UseDesc = rsp.Entry.UseTs.String()
	return rsp
}

type linkTraceWarp struct {
	lst       collections.List[response.LinkTraceVO]
	rgbaIndex int                                  // 实现不同服务的调用，用颜色区分。这里通过服务入口调用时，使用索引（对应数组颜色）
	lstPO     collections.List[trace.TraceContext] // 数据库读的集合
	startTs   int64                                // 初始开始时间（微秒）
	endTs     int64                                // 初始结束时间（微秒）
	TotalUse  float64                              // 总共时间（微秒）
	PreDetail trace.TraceDetail                    // 上一次的执行明细。用来解决两个服务间时间不同步（服务器的时间没有同步）
}

// 服务调用入口
func (receiver *linkTraceWarp) addEntry(po trace.TraceContext) {
	receiver.lstPO.RemoveAll(func(item trace.TraceContext) bool {
		return item.AppName == po.AppName && item.ParentAppName == po.ParentAppName && item.TraceType == po.TraceType && (item.TraceLevel == po.TraceLevel) && item.StartTs == po.StartTs && item.EndTs == po.EndTs
	})

	receiver.rgbaIndex++
	// 添加服务的入口
	entryTrace := response.LinkTraceVO{
		Rgba: response.RgbaList[receiver.rgbaIndex], AppId: po.AppId, AppIp: po.AppIp, AppName: po.AppName, UseTs: float64(po.UseTs.Microseconds()), UseDesc: po.UseTs.String(),
		StartTs:   float64(po.StartTs - receiver.startTs),
		Exception: po.Exception,
	}
	entryTrace.StartRate = entryTrace.StartTs / receiver.TotalUse * 100
	entryTrace.UseRate = entryTrace.UseTs / receiver.TotalUse * 100
	// 通常说明不同服务间的机器时间不同步
	if entryTrace.StartRate > 100 || entryTrace.StartRate < 0 {
		// 使用上一个入口的结束时间
		entryTrace.StartRate = float64(receiver.PreDetail.EndTs-receiver.startTs) / receiver.TotalUse * 100
	}
	switch po.TraceType {
	case eumTraceType.WebApi:
		entryTrace.Caption = fmt.Sprintf("收到%s请求【%s】 => %s", po.WebRequestIp, po.WebMethod, po.WebPath)
		entryTrace.Desc = fmt.Sprintf("%s 客户端IP：%s", po.WebContentType, po.WebRequestIp)
	case eumTraceType.MqConsumer:
		entryTrace.Caption = fmt.Sprintf("MQ订阅 => %s %s %s", po.ConsumerServer, po.ConsumerQueueName, po.ConsumerRoutingKey)
	case eumTraceType.QueueConsumer:
		entryTrace.Caption = fmt.Sprintf("本地Queue订阅 => %s", po.ConsumerQueueName)
	case eumTraceType.EventConsumer:
		entryTrace.Caption = fmt.Sprintf("事件订阅 => %s %s", po.ConsumerServer, po.ConsumerQueueName)
	case eumTraceType.FSchedule:
		entryTrace.Caption = fmt.Sprintf("任务调度 => 任务组：%s 任务ID：%v", po.TaskGroupName, po.TaskId)
		dataJson, _ := po.TaskData.MarshalJSON()
		entryTrace.Desc = fmt.Sprintf("参数 %s", string(dataJson))
	case eumTraceType.Task:
		entryTrace.Caption = fmt.Sprintf("本地任务 => %s", po.TaskName)
	case eumTraceType.WatchKey:
		entryTrace.Caption = fmt.Sprintf("监控KEY => %s", po.WatchKey)
	}
	receiver.lst.Add(entryTrace)
	receiver.addDetail(po)
	receiver.rgbaIndex--
}

// 服务所属的明细
func (receiver *linkTraceWarp) addDetail(po trace.TraceContext) {
	po.List.Foreach(func(traceDetail **trace.TraceDetail) {
		detail := *traceDetail
		useTs := time.Duration(detail.EndTs-detail.StartTs) * time.Microsecond
		detailTrace := response.LinkTraceVO{
			Rgba: response.RgbaList[receiver.rgbaIndex], AppId: po.AppId, AppIp: po.AppIp, AppName: po.AppName,
			StartTs: float64(detail.StartTs - receiver.startTs),
			UseTs:   float64(useTs.Microseconds()), UseDesc: useTs.String(),
			Exception: detail.Exception,
		}

		detailTrace.StartRate = detailTrace.StartTs / receiver.TotalUse * 100
		detailTrace.UseRate = detailTrace.UseTs / receiver.TotalUse * 100
		// 通常说明不同服务间的机器时间不同步
		if detailTrace.StartRate > 100 {
			// 使用上一个入口的结束时间
			detailTrace.StartRate = float64(receiver.PreDetail.EndTs-receiver.startTs) / receiver.TotalUse * 100
		}
		switch detail.CallType {
		case eumCallType.Database:
			if detail.DbTableName == "" && detail.DbSql == "" {
				detailTrace.Caption = fmt.Sprintf("%s[打开数据库] => %s %s", detail.Comment, detail.DbName, detail.DbConnectionString)
			} else {
				if len(detail.DbSql) < 400 {
					detail.DbSql = strings.ReplaceAll(detail.DbSql, "\n", "")
					detailTrace.Caption = fmt.Sprintf("SQL <span class=\"el-tag el-tag--danger el-tag--small el-tag--light\">%s</span> => <span style='background-color: #ead996;'>%s</span> 影响%v行", detail.Comment, detail.DbSql, detail.DbRowsAffected)
				} else {
					detailTrace.Caption = fmt.Sprintf("SQL <span class=\"el-tag el-tag--danger el-tag--small el-tag--light\">%s</span> => %s.<b>%s</b> 影响%v行", detail.Comment, detail.DbName, detail.DbTableName, detail.DbRowsAffected)
				}
			}
			detailTrace.Desc = detail.DbSql
		case eumCallType.Http:
			detailTrace.Caption = fmt.Sprintf("调用http <span class=\"el-tag el-tag--danger el-tag--small el-tag--light\">%s</span> => %v %s <span style='background-color: #ead996;'>%s</span>", detail.Comment, detail.HttpStatusCode, detail.HttpMethod, detail.HttpUrl)
			lstHeader := collections.NewList[string]()
			for k, v := range detail.HttpHeaders.ToMap() {
				lstHeader.Add(fmt.Sprintf("%s=%v", k, v))
			}
			detailTrace.Desc = fmt.Sprintf("头部：%s 入参：%s 出参：%s", lstHeader.ToString(","), detail.HttpRequestBody, detail.HttpResponseBody)
		case eumCallType.Grpc:
			detailTrace.Caption = fmt.Sprintf("调用http <span class=\"el-tag el-tag--danger el-tag--small el-tag--light\">%s</span> => %v %s <span style='background-color: #ead996;'>%s</span>", detail.Comment, detail.GrpcStatusCode, detail.GrpcMethod, detail.GrpcUrl)
			lstHeader := collections.NewList[string]()
			for k, v := range detail.GrpcHeaders.ToMap() {
				lstHeader.Add(fmt.Sprintf("%s=%v", k, v))
			}
			detailTrace.Desc = fmt.Sprintf("头部：%s 入参：%s 出参：%s", lstHeader.ToString(","), detail.GrpcRequestBody, detail.GrpcResponseBody)
		case eumCallType.Redis:
			detailTrace.Caption = fmt.Sprintf("Redis <span class=\"el-tag el-tag--danger el-tag--small el-tag--light\">%s</span> => <span style='background-color: #ead996;'>%s</span> %s %s 影响%v行", detail.Comment, detail.MethodName, detail.RedisKey, detail.RedisField, detail.RedisRowsAffected)
			detailTrace.Desc = fmt.Sprintf("%s %s", detail.RedisKey, detail.RedisField)
		case eumCallType.Mq:
			if detail.MethodName == "Send" {
				detailTrace.Caption = fmt.Sprintf("MQ发消息 <span class=\"el-tag el-tag--danger el-tag--small el-tag--light\">%s</span> => %s <span style='background-color: #ead996;'>%s</span> %s", detail.Comment, detail.MqServer, detail.MqExchange, detail.MqRoutingKey)
			} else {
				detailTrace.Caption = fmt.Sprintf("MQ <span class=\"el-tag el-tag--danger el-tag--small el-tag--light\">%s</span> %s => %s <span style='background-color: #ead996;'>%s</span> %s", detail.Comment, detail.MethodName, detail.MqServer, detail.MqExchange, detail.MqRoutingKey)
			}
			detailTrace.Desc = fmt.Sprintf("%s %s %s", detail.MqServer, detail.MqExchange, detail.MqRoutingKey)
		case eumCallType.Elasticsearch:
			detailTrace.Caption = fmt.Sprintf("ES <span class=\"el-tag el-tag--danger el-tag--small el-tag--light\">%s</span> => %s %s", detail.Comment, detail.EsIndexName, detail.EsAliasesName)
			detailTrace.Desc = fmt.Sprintf("%s %s", detail.EsIndexName, detail.EsAliasesName)
		case eumCallType.Etcd:
			detailTrace.Caption = fmt.Sprintf("Etcd <span class=\"el-tag el-tag--danger el-tag--small el-tag--light\">%s</span> => %s %v", detail.Comment, detail.EtcdKey, detail.EtcdLeaseID)
			detailTrace.Desc = fmt.Sprintf("%s %v", detail.EtcdKey, detail.EtcdLeaseID)
		case eumCallType.Hand:
			detailTrace.Caption = fmt.Sprintf("<span class=\"el-tag el-tag--danger el-tag--small el-tag--light\">%s %s</span>", detail.Comment, detail.HandName)
			detailTrace.Desc = detail.HandName
		case eumCallType.EventPublish:
			detailTrace.Caption = fmt.Sprintf("事件订阅 <span class=\"el-tag el-tag--danger el-tag--small el-tag--light\">%s</span> => <span style='background-color: #ead996;'>%s</span>", detail.Comment, detail.EventName)
			detailTrace.Desc = detail.EventName
		}

		detailTrace.Caption = strings.ReplaceAll(detailTrace.Caption, "<span class=\"el-tag el-tag--danger el-tag--small el-tag--light\"></span>", "")
		receiver.lst.Add(detailTrace)

		// 在明细执行期间，会穿插下游服务。所以通过查找的方式来获取下游。然后在回到当前明细
		// a --> b -- > a  --> c -- b
		var nextEntry trace.TraceContext
		switch detail.CallType {
		case eumCallType.Http:
			// 查找串联的服务
			nextEntry = receiver.lstPO.Where(func(item trace.TraceContext) bool {
				return item.ParentAppName == detailTrace.AppName && item.TraceType == eumTraceType.WebApi && (item.TraceLevel == po.TraceLevel+1)
			}).OrderBy(func(item trace.TraceContext) any {
				return item.StartTs
			}).First()
		case eumCallType.EventPublish:
			nextEntry = receiver.lstPO.Where(func(item trace.TraceContext) bool {
				return item.ParentAppName == detailTrace.AppName && item.TraceType == eumTraceType.EventConsumer && (item.TraceLevel == po.TraceLevel+1)
			}).OrderBy(func(item trace.TraceContext) any {
				return item.StartTs
			}).First()
		case eumCallType.Mq:
			nextEntry = receiver.lstPO.Where(func(item trace.TraceContext) bool {
				return item.ParentAppName == detailTrace.AppName && item.TraceType == eumTraceType.MqConsumer && (item.TraceLevel == po.TraceLevel+1)
			}).OrderBy(func(item trace.TraceContext) any {
				return item.StartTs
			}).First()
		}
		if nextEntry.TraceId != "" {
			receiver.PreDetail = *detail
			receiver.addEntry(nextEntry)
		}
	})
}
