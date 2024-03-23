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
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/context"
)

var route = []webapi.Route{
	{"POST", "/apps/add", appsApp.Add, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
	{"POST", "/apps/update", appsApp.Update, "", []context.IFilter{application.Jwt{}}, []string{"req", "", ""}},
	{"POST", "/apps/delete", appsApp.Delete, "", []context.IFilter{application.Jwt{}}, []string{"appName", "", ""}},
	{"POST", "/apps/list", appsApp.List, "", []context.IFilter{application.Jwt{}}, []string{"clusterId", "", ""}},
	{"POST", "/apps/info", appsApp.Info, "", []context.IFilter{application.Jwt{}}, []string{"appName", ""}},
	{"POST", "/apps/build/add", appsApp.BuildAdd, "", []context.IFilter{application.Jwt{}}, []string{"appName", "clusterId", "", ""}},
	{"POST", "/apps/build/list", appsApp.BuildList, "", []context.IFilter{application.Jwt{}}, []string{"appName", "pageSize", "pageIndex", ""}},
	{"GET", "/apps/build/view-{buildId}", appsApp.View, "", []context.IFilter{}, []string{"buildId"}},
	{"POST", "/apps/build/stop", appsApp.Stop, "", []context.IFilter{application.Jwt{}}, []string{""}},
	{"POST", "/apps/build/syncDockerImage", appsApp.SyncDockerImage, "", []context.IFilter{application.Jwt{}}, []string{"clusterId", "appName", "", "", ""}},
	{"POST", "/apps/deleteService", appsApp.DeleteService, "", []context.IFilter{application.Jwt{}}, []string{"appName", "", ""}},
	{"POST", "/apps/updateDockerImage", appsApp.UpdateDockerImage, "", []context.IFilter{}, []string{"appName", "dockerImage", "buildNumber", "clusterId", "dockerHub", "dockerUserName", "dockerUserPwd", "", "", "", ""}},
	{"POST", "/apps/build/clearDockerImage", appsApp.ClearDockerImage, "", []context.IFilter{application.Jwt{}}, []string{""}},
	{"POST", "/apps/build/restartDocker", appsApp.RestartDocker, "", []context.IFilter{application.Jwt{}}, []string{"clusterId", "appName", "", ""}},
	{"POST", "/apps/register", appsApp.Register, "", []context.IFilter{}, []string{"req", ""}},
	{"POST", "/cluster/add", clusterApp.Add, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
	{"POST", "/cluster/update", clusterApp.Update, "", []context.IFilter{application.Jwt{}}, []string{"req", ""}},
	{"POST", "/cluster/list", clusterApp.List, "", []context.IFilter{application.Jwt{}}, []string{""}},
	{"POST", "/cluster/delete", clusterApp.Delete, "", []context.IFilter{application.Jwt{}}, []string{"clusterId", ""}},
	{"POST", "/configure/list", configureApp.List, "", []context.IFilter{}, []string{"appName", ""}},
	{"GET", "/flog/list", flogApp.List, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "logContent", "logLevel", "pageSize", "pageIndex", ""}},
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
	{"GET", "/linkTrace/taskList", linkTraceApp.TaskList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "taskName", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
	{"GET", "/linkTrace/fScheduleList", linkTraceApp.FScheduleList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "taskName", "taskGroupId", "taskId", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
	{"GET", "/linkTrace/consumerList", linkTraceApp.ConsumerList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "server", "queueName", "routingKey", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
	{"GET", "/linkTrace/slowDbList", linkTraceApp.SlowDbList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "dbName", "tableName", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
	{"GET", "/linkTrace/slowEsList", linkTraceApp.SlowEsList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "indexName", "aliasesName", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
	{"GET", "/linkTrace/slowEtcdList", linkTraceApp.SlowEtcdList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "key", "leaseID", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
	{"GET", "/linkTrace/slowHandList", linkTraceApp.SlowHandList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "name", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
	{"GET", "/linkTrace/slowHttpList", linkTraceApp.SlowHttpList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "method", "url", "requestBody", "responseBody", "statusCode", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
	{"GET", "/linkTrace/slowMqList", linkTraceApp.SlowMqList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "server", "exchange", "routingKey", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
	{"GET", "/linkTrace/slowRedisList", linkTraceApp.SlowRedisList, "", []context.IFilter{application.Jwt{}}, []string{"traceId", "appName", "appIp", "key", "field", "searchUseTs", "onlyViewException", "startMin", "pageSize", "pageIndex", ""}},
	{"POST", "/linkTrace/upload", linkTraceApp.Upload, "", []context.IFilter{}, []string{"req", ""}},
	{"POST", "/user/passport/Login", login.Login, "", []context.IFilter{}, []string{"req", ""}},
}
