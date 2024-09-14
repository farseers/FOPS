package context

import (
	"fops/domain/linkTrace"
	"fops/infrastructure/repository/model"
	"github.com/farseer-go/data"
)

var CHContext *chContext

// EsContext 链路追踪上下文
type chContext struct {
	data.IInternalContext
	TraceContextView    data.TableSet[model.TraceContextViewPO]    `data:"name=v_link_trace;migrate"`
	TraceContext        data.TableSet[model.TraceContextPO]        `data:"name=link_trace;migrate=ReplacingMergeTree() ORDER BY (trace_type,app_name,parent_app_name,app_ip,trace_id,start_ts) PARTITION BY (trace_type)"`
	TraceDetailDatabase data.TableSet[model.TraceDetailDatabasePO] `data:"name=trace_detail_database;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,db_name,table_name,trace_id,start_ts) PARTITION BY (app_name)"`
	TraceDetailEs       data.TableSet[model.TraceDetailEsPO]       `data:"name=trace_detail_es;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,index_name,aliases_name,trace_id,start_ts) PARTITION BY (app_name)"`
	TraceDetailEtcd     data.TableSet[model.TraceDetailEtcdPO]     `data:"name=trace_detail_etcd;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,key,trace_id,start_ts) PARTITION BY (app_name)"`
	TraceDetailHand     data.TableSet[model.TraceDetailHandPO]     `data:"name=trace_detail_hand;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,name,trace_id,start_ts) PARTITION BY (app_name)"`
	TraceDetailHttp     data.TableSet[model.TraceDetailHttpPO]     `data:"name=trace_detail_http;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,method,url,trace_id,start_ts) PARTITION BY (app_name)"`
	TraceDetailGrpc     data.TableSet[model.TraceDetailGrpcPO]     `data:"name=trace_detail_grpc;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,method,url,trace_id,start_ts) PARTITION BY (app_name)"`
	TraceDetailMq       data.TableSet[model.TraceDetailMqPO]       `data:"name=trace_detail_mq;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,server,exchange,routing_key,trace_id,start_ts) PARTITION BY (app_name)"`
	TraceDetailRedis    data.TableSet[model.TraceDetailRedisPO]    `data:"name=trace_detail_redis;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,key,field,trace_id,start_ts) PARTITION BY (app_name)"`
	LogData             data.TableSet[model.LogDataPO]             `data:"name=log_data;migrate=ReplacingMergeTree() ORDER BY (app_name,component,log_level,app_ip,trace_id,log_id,create_at) PARTITION BY (app_name)"`
	Visits              data.TableSet[model.VisitsPO]              `data:"name=visits;migrate=ReplacingMergeTree() ORDER BY (create_at,app_name,visits_node_prefix,visits_node,trace_type) PARTITION BY toYYYYMM(create_at)"`
	MonitorData         data.TableSet[model.MonitorDataPO]         `data:"name=monitor_data;migrate=ReplacingMergeTree() ORDER BY (app_name,create_at) PARTITION BY (app_name)"`
}

// InitChContextContext 初始化上下文
func InitChContextContext() {
	if linkTrace.Config.ConnString == "" {
		panic("[farseer.yaml]FOPS.LinkTrace.ConnString，配置不正确")
	}
	data.RegisterInternalContext("LinkTrace", linkTrace.Config.ConnString)
	CHContext = data.NewContext[chContext]("LinkTrace")
}
