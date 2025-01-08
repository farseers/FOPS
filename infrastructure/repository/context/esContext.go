package context

import (
	"fops/domain/linkTrace"
	"fops/infrastructure/repository/model"
	"github.com/farseer-go/elasticSearch"
)

var ESContext *esContext

// EsContext 链路追踪上下文
type esContext struct {
	TraceContext        elasticSearch.IndexSet[model.TraceContextPO]        `es:"index=link_trace_yyyy_MM;alias=link_trace;shards=1;replicas=0;refresh=3"`
	TraceDetailDatabase elasticSearch.IndexSet[model.TraceDetailDatabasePO] `es:"index=trace_detail_database_yyyy_MM;alias=trace_detail_database;shards=1;replicas=0;refresh=3"`
	TraceDetailEs       elasticSearch.IndexSet[model.TraceDetailEsPO]       `es:"index=trace_detail_es_yyyy_MM;alias=trace_detail_es;shards=1;replicas=0;refresh=3"`
	TraceDetailEtcd     elasticSearch.IndexSet[model.TraceDetailEtcdPO]     `es:"index=trace_detail_etcd_yyyy_MM;alias=trace_detail_etcd;shards=1;replicas=0;refresh=3"`
	TraceDetailHand     elasticSearch.IndexSet[model.TraceDetailHandPO]     `es:"index=trace_detail_hand_yyyy_MM;alias=trace_detail_hand;shards=1;replicas=0;refresh=3"`
	TraceDetailHttp     elasticSearch.IndexSet[model.TraceDetailHttpPO]     `es:"index=trace_detail_http_yyyy_MM;alias=trace_detail_http;shards=1;replicas=0;refresh=3"`
	TraceDetailGrpc     elasticSearch.IndexSet[model.TraceDetailGrpcPO]     `es:"index=trace_detail_grpc_yyyy_MM;alias=trace_detail_http;shards=1;replicas=0;refresh=3"`
	TraceDetailMq       elasticSearch.IndexSet[model.TraceDetailMqPO]       `es:"index=trace_detail_mq_yyyy_MM;alias=trace_detail_mq;shards=1;replicas=0;refresh=3"`
	TraceDetailRedis    elasticSearch.IndexSet[model.TraceDetailRedisPO]    `es:"index=trace_detail_redis_yyyy_MM;alias=trace_detail_redis;shards=1;replicas=0;refresh=3"`
	LogData             elasticSearch.IndexSet[model.LogDataPO]             `es:"index=log_data_yyyy_MM;alias=log_data;shards=1;replicas=0;refresh=3"`
}

// initEsContext 初始化上下文
func initEsContext() {
	if linkTrace.Config.ConnString == "" {
		panic("[farseer.yaml]FOPS.LinkTrace.ConnString，配置不正确")
	}
	elasticSearch.RegisterInternalContext("LinkTrace", linkTrace.Config.ConnString)
	ESContext = elasticSearch.NewContext[esContext]("LinkTrace")
}
