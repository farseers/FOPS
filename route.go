// 该文件由fsctl route命令自动生成，请不要手动修改此文件
package main

import (
	"fops/application"
	"fops/application/appsApp"
	"fops/application/clusterApp"
	"fops/application/configureApp"
	"fops/application/flogApp"
	"fops/application/gitApp"
	"fops/application/linkTraceApp"
	"fops/application/login"
	"fops/application/monitorApp"
	"fops/application/terminalApp"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/context"
)

var route = []webapi.Route{
    {"POST", "/apps/add", appsApp.Add, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/apps/update", appsApp.Update, "", []context.IFilter{application.Jwt{}}, []string{"req", "", ""}},
    {"POST", "/apps/delete", appsApp.Delete, "", []context.IFilter{application.Jwt{}}, []string{"appName", ""}},
    {"POST", "/apps/dropDownList", appsApp.DropDownList, "", []context.IFilter{application.Jwt{}}, []string{"isAll", "", ""}},
    {"POST", "/apps/list", appsApp.List, "", []context.IFilter{application.Jwt{}}, []string{"isSys", "", "", "", ""}},
    {"POST", "/apps/info", appsApp.Info, "", []context.IFilter{application.Jwt{}}, []string{"appName", "", ""}},
    {"POST", "/apps/syncWorkflows", appsApp.SyncWorkflows, "", []context.IFilter{application.Jwt{}}, []string{"appName", "", ""}},
    {"POST", "/apps/build/add", appsApp.BuildAdd, "", []context.IFilter{}, []string{"appName", "workflowsName", "branchName", "", "", ""}},
    {"POST", "/apps/build/list", appsApp.BuildList, "", []context.IFilter{application.Jwt{}}, []string{"appName", "pageSize", "pageIndex", ""}},
    {"GET", "/apps/build/view-{buildId}", appsApp.View, "", []context.IFilter{}, []string{"buildId"}},
    {"POST", "/apps/build/stop", appsApp.Stop, "", []context.IFilter{application.Jwt{}}, []string{"buildId", ""}},
    {"POST", "/apps/updateDockerImage", appsApp.UpdateDockerImage, "", []context.IFilter{}, []string{"appName", "dockerImage", "updateDelay", "buildNumber", "dockerHub", "dockerUserName", "dockerUserPwd", "", ""}},
    {"POST", "/apps/build/clearDockerImage", appsApp.ClearDockerImage, "", []context.IFilter{application.Jwt{}}, []string{}},
    {"POST", "/apps/restartDocker", appsApp.RestartDocker, "", []context.IFilter{application.Jwt{}}, []string{"appName", "", ""}},
    {"POST", "/apps/setReplicas", appsApp.SetReplicas, "", []context.IFilter{application.Jwt{}}, []string{"appName", "dockerReplicas", ""}},
    {"POST", "/apps/deleteService", appsApp.DeleteService, "", []context.IFilter{application.Jwt{}}, []string{"appName", ""}},
    {"POST", "/apps/logs/dockerSwarm", appsApp.DockerSwarm, "", []context.IFilter{application.Jwt{}}, []string{"appName", "tailCount"}},
    {"POST", "/cluster/add", clusterApp.Add, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/cluster/update", clusterApp.Update, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/cluster/list", clusterApp.List, "", []context.IFilter{application.Jwt{}}, []string{""}},
    {"POST", "/cluster/delete", clusterApp.Delete, "", []context.IFilter{application.Jwt{}}, []string{"clusterId", ""}},
    {"GET", "/cluster/nodeList", clusterApp.NodeList, "", []context.IFilter{application.Jwt{}}, []string{""}},
    {"POST", "/configure/list", configureApp.List, "", []context.IFilter{}, []string{"appName", ""}},
    {"GET", "/configure/allList", configureApp.AllList, "", []context.IFilter{application.Jwt{}}, []string{""}},
    {"POST", "/configure/add", configureApp.Add, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/configure/update", configureApp.Update, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/configure/rollback", configureApp.Rollback, "", []context.IFilter{application.Jwt{}}, []string{"appName", "key", ""}},
    {"POST", "/configure/delete", configureApp.Delete, "", []context.IFilter{application.Jwt{}}, []string{"appName", "key", ""}},
    {"GET", "/flog/list", flogApp.List, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "logContent", "minute", "logLevel", "pageSize", "pageIndex", ""}},
    {"POST", "/flog/delete", flogApp.Delete, "", []context.IFilter{application.Jwt{}}, []string{""}},
    {"GET", "/flog/info-{id}", flogApp.Info, "", []context.IFilter{application.Jwt{}}, []string{"id", ""}},
    {"GET", "/flog/StatCount", flogApp.StatCount, "", []context.IFilter{application.Jwt{}}, []string{"appName", ""}},
    {"POST", "/flog/upload", flogApp.Upload, "", []context.IFilter{}, []string{"req", ""}},
    {"POST", "/git/add", gitApp.Add, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/git/update", gitApp.Update, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/git/list", gitApp.List, "", []context.IFilter{application.Jwt{}}, []string{"isApp", ""}},
    {"POST", "/git/delete", gitApp.Delete, "", []context.IFilter{application.Jwt{}}, []string{"gitId", ""}},
    {"POST", "/git/info", gitApp.Info, "", []context.IFilter{application.Jwt{}}, []string{"gitId", ""}},
    {"GET", "/linkTrace/info/{traceId}", linkTraceApp.Info, "", []context.IFilter{}, []string{"traceId", ""}},
    {"GET", "/linkTrace/webApiList", linkTraceApp.WebApiList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "requestIp", "searchUrl", "statusCode", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"POST", "/linkTrace/delete", linkTraceApp.Delete, "", []context.IFilter{application.Jwt{}}, []string{"traceType", ""}},
    {"GET", "/linkTrace/webSocketList", linkTraceApp.WebSocketList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "requestIp", "searchUrl", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"GET", "/linkTrace/taskList", linkTraceApp.TaskList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "taskName", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"GET", "/linkTrace/fScheduleList", linkTraceApp.FScheduleList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "taskName", "taskGroupId", "taskId", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"GET", "/linkTrace/consumerList", linkTraceApp.ConsumerList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "server", "queueName", "routingKey", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"GET", "/linkTrace/queueList", linkTraceApp.QueueList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "queueName", "routingKey", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"GET", "/linkTrace/eventList", linkTraceApp.EventList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "queueName", "routingKey", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"POST", "/linkTrace/deleteSlow", linkTraceApp.DeleteSlow, "", []context.IFilter{application.Jwt{}}, []string{"dbName", ""}},
    {"GET", "/linkTrace/slowDbList", linkTraceApp.SlowDbList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "dbName", "tableName", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"GET", "/linkTrace/slowEsList", linkTraceApp.SlowEsList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "indexName", "aliasesName", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"GET", "/linkTrace/slowEtcdList", linkTraceApp.SlowEtcdList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "key", "leaseID", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"GET", "/linkTrace/slowHandList", linkTraceApp.SlowHandList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "name", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"GET", "/linkTrace/slowHttpList", linkTraceApp.SlowHttpList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "method", "url", "body", "statusCode", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"GET", "/linkTrace/slowMqList", linkTraceApp.SlowMqList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "server", "exchange", "routingKey", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"GET", "/linkTrace/slowRedisList", linkTraceApp.SlowRedisList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "key", "field", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
    {"POST", "/linkTrace/upload", linkTraceApp.Upload, "", []context.IFilter{}, []string{"req", ""}},
    {"GET", "/linkTrace/visits", linkTraceApp.Visits, "", []context.IFilter{}, []string{"request", ""}},
    {"POST", "/user/passport/Login", login.Login, "", []context.IFilter{}, []string{"req", ""}},
    {"POST", "/user/passport/changePwd", login.ChangePwd, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/monitor/appList", monitorApp.DropDownListAppInfo, "", []context.IFilter{application.Jwt{}}, []string{""}},
    {"POST", "/monitor/ruleList", monitorApp.ToListPageRule, "", []context.IFilter{application.Jwt{}}, []string{"appName", "pageSize", "pageIndex", ""}},
    {"POST", "/monitor/delRule", monitorApp.DeleteRule, "", []context.IFilter{application.Jwt{}}, []string{"id", ""}},
    {"POST", "/monitor/infoRule", monitorApp.ToEntityRule, "", []context.IFilter{application.Jwt{}}, []string{"id", ""}},
    {"POST", "/monitor/saveRule", monitorApp.SaveRule, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/monitor/noticeList", monitorApp.ToListPageNotice, "", []context.IFilter{application.Jwt{}}, []string{"name", "pageSize", "pageIndex", ""}},
    {"POST", "/monitor/delNotice", monitorApp.DeleteNotice, "", []context.IFilter{application.Jwt{}}, []string{"id", ""}},
    {"POST", "/monitor/infoNotice", monitorApp.ToEntityNotice, "", []context.IFilter{application.Jwt{}}, []string{"id", ""}},
    {"POST", "/monitor/saveNotice", monitorApp.SaveNotice, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/monitor/dataList", monitorApp.ToListPageData, "", []context.IFilter{application.Jwt{}}, []string{"appName", "pageSize", "pageIndex", ""}},
    {"POST", "/monitor/noticeLogList", monitorApp.ToListPageNoticeLog, "", []context.IFilter{application.Jwt{}}, []string{"appName", "pageSize", "pageIndex", ""}},
    {"POST", "/monitor/noticeLogNoReadList", monitorApp.ToListPageNoticeLogNoRead, "", []context.IFilter{application.Jwt{}}, []string{"top", ""}},
    {"POST", "/monitor/allRead", monitorApp.UpdateNoticeLogRead, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/monitor/delNoticeLog", monitorApp.DeleteNoticeLog, "", []context.IFilter{application.Jwt{}}, []string{""}},
    {"POST", "/monitor/drpBaseList", monitorApp.DrpBaseList, "", []context.IFilter{application.Jwt{}}, []string{"baseType"}},
    {"WS", "/ws/monitor", monitorApp.WsReceive, "", []context.IFilter{}, []string{"context", ""}},
    {"WS", "/ws/notice", monitorApp.WsNotice, "", []context.IFilter{}, []string{"context", ""}},
    {"POST", "/terminal/clientList", terminalApp.ClientList, "", []context.IFilter{application.Jwt{}}, []string{"pageSize", "pageIndex", ""}},
    {"POST", "/terminal/clientAdd", terminalApp.ClientAdd, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/terminal/clientUpdate", terminalApp.ClientUpdate, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
    {"POST", "/terminal/clientDel", terminalApp.ClientDel, "", []context.IFilter{application.Jwt{}}, []string{"id", ""}},
    {"POST", "/terminal/clientInfo", terminalApp.ClientInfo, "", []context.IFilter{application.Jwt{}}, []string{"id", ""}},
    {"WS", "/terminal/ws/ssh", terminalApp.WsSsh, "", []context.IFilter{application.Jwt{}}, []string{"context", ""}},
    {"WS", "/terminal/ws/sshByLogin", terminalApp.WsSshByLogin, "", []context.IFilter{application.Jwt{}}, []string{"context", ""}},
}
