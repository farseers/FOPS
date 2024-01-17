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
	TraceContext        data.TableSet[model.TraceContextPO]        `data:"name=link_trace;migrate=ReplacingMergeTree() ORDER BY (trace_level,app_name,parent_app_name,app_ip,app_id,trace_id,start_ts) PARTITION BY toYYYYMM(create_at)"`
	TraceDetailDatabase data.TableSet[model.TraceDetailDatabasePO] `data:"name=trace_detail_database;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,db_name,table_name,start_ts) PARTITION BY toYYYYMM(create_at)"`
	TraceDetailEs       data.TableSet[model.TraceDetailEsPO]       `data:"name=trace_detail_es;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,index_name,aliases_name,start_ts) PARTITION BY toYYYYMM(create_at)"`
	TraceDetailEtcd     data.TableSet[model.TraceDetailEtcdPO]     `data:"name=trace_detail_etcd;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,key,start_ts) PARTITION BY toYYYYMM(create_at)"`
	TraceDetailHand     data.TableSet[model.TraceDetailHandPO]     `data:"name=trace_detail_hand;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,name,start_ts) PARTITION BY toYYYYMM(create_at)"`
	TraceDetailHttp     data.TableSet[model.TraceDetailHttpPO]     `data:"name=trace_detail_http;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,method,url,start_ts) PARTITION BY toYYYYMM(create_at)"`
	TraceDetailGrpc     data.TableSet[model.TraceDetailGrpcPO]     `data:"name=trace_detail_grpc;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,method,url,start_ts) PARTITION BY toYYYYMM(create_at)"`
	TraceDetailMq       data.TableSet[model.TraceDetailMqPO]       `data:"name=trace_detail_mq;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,server,exchange,routing_key,start_ts) PARTITION BY toYYYYMM(create_at)"`
	TraceDetailRedis    data.TableSet[model.TraceDetailRedisPO]    `data:"name=trace_detail_redis;migrate=ReplacingMergeTree() ORDER BY (app_name,parent_app_name,app_ip,app_id,trace_id,key,field,start_ts) PARTITION BY toYYYYMM(create_at)"`
	LogData             data.TableSet[model.LogDataPO]             `data:"name=log_data;migrate=ReplacingMergeTree() ORDER BY (app_name,component,log_level,app_ip,app_id,trace_id,create_at,log_id) PARTITION BY toYYYYMM(create_at)"`
}

// InitChContextContext 初始化上下文
func InitChContextContext() {
	if linkTrace.Config.ConnString == "" {
		panic("[farseer.yaml]FOPS.LinkTrace.ConnString，配置不正确")
	}
	data.RegisterInternalContext("LinkTrace", linkTrace.Config.ConnString)
	CHContext = data.NewContext[chContext]("LinkTrace")
}
