// @area /apps/
package appsApp

import (
	"fops/application/appsApp/request"
	"fops/application/appsApp/response"
	"fops/domain/apps"
	"fops/domain/cluster"
	"fops/domain/logData"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/utils/file"
	"path/filepath"
	"strings"
)

// Add 添加应用
// @post add
// @filter application.Jwt
func Add(req request.AddRequest, appsRepository apps.Repository) {
	do := mapper.Single[apps.DomainObject](req)
	exception.ThrowWebExceptionBool(appsRepository.IsExists(req.AppName), 403, "应用不能重复")
	// 删除末尾的/
	if strings.HasSuffix(do.AdditionalScripts, "\\") {
		do.AdditionalScripts = do.AdditionalScripts[:len(do.AdditionalScripts)-1]
	}

	// 添加
	err := appsRepository.Add(do)
	exception.ThrowWebExceptionError(403, err)
}

// Update 修改应用
// @post update
// @filter application.Jwt
func Update(req request.UpdateRequest, appsRepository apps.Repository, appsIDockerSwarmDevice apps.IDockerSwarmDevice) {
	do := appsRepository.ToEntity(req.AppName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	// 判断副本数量是否有变更
	if do.DockerReplicas != req.DockerReplicas {
		c := make(chan string, 100)
		if !appsIDockerSwarmDevice.SetReplicas(cluster.DomainObject{}, req.AppName, req.DockerReplicas, c) {
			lstLog := collections.NewListFromChan(c)
			exception.ThrowWebExceptionf(403, "更新副本失败:<br />%s", lstLog.ToString("<br />"))
		}
	}

	// 更新应用信息
	do = mapper.Single[apps.DomainObject](req)
	// 删除末尾的/
	if strings.HasSuffix(do.AdditionalScripts, "\\") {
		do.AdditionalScripts = do.AdditionalScripts[:len(do.AdditionalScripts)-1]
	}
	err := appsRepository.UpdateApp(do)
	exception.ThrowWebExceptionError(403, err)
}

// Delete 删除应用
// @post delete
// @filter application.Jwt
func Delete(appName string, appsRepository apps.Repository, appsIDockerSwarmDevice apps.IDockerSwarmDevice) {
	exception.ThrowWebExceptionBool(strings.Trim(appName, "") == "", 403, "参数不完整")
	// 删除服务
	c := make(chan string, 100)
	appsIDockerSwarmDevice.DeleteService(appName, c)

	// 删除应用
	_, err := appsRepository.Delete(appName)
	exception.ThrowWebExceptionError(403, err)
}

// List 应用列表
// @post list
// @filter application.Jwt
func List(clusterId int64, appsRepository apps.Repository, logDataRepository logData.Repository) collections.List[response.AppsResponse] {
	lstDO := appsRepository.ToList()
	lstGit := appsRepository.ToGitListAll(-1)
	countList := logDataRepository.StatCount()

	lst := collections.NewList[response.AppsResponse]()
	lstDO.Foreach(func(item *apps.DomainObject) {
		appsResponse := doToAppsResponse(clusterId, *item)
		appsResponse.AppGitName = lstGit.Where(func(gitItem apps.GitEO) bool {
			return item.AppGit == parse.ToInt64(gitItem.Id)
		}).First().Name

		appsResponse.LogErrorCount = countList.Where(func(logItem logData.LogCountEO) bool {
			return item.AppName == logItem.AppName && logItem.LogLevel == eumLogLevel.Error
		}).First().LogCount

		appsResponse.LogWaringCount = countList.Where(func(logItem logData.LogCountEO) bool {
			return item.AppName == logItem.AppName && logItem.LogLevel == eumLogLevel.Warning
		}).First().LogCount

		// 获取工作流文件名称
		appsResponse.WorkflowsNames = file.GetFiles(item.GetWorkflowsDir(), "*.yml", true)
		for i, v := range appsResponse.WorkflowsNames {
			v = filepath.Base(v)
			appsResponse.WorkflowsNames[i] = v[:strings.Index(v, ".yml")]
		}
		lst.Add(appsResponse)
	})
	return lst
}

// Info 查询应用
// @post info
// @filter application.Jwt
func Info(appName string, appsRepository apps.Repository) response.AppsResponse {
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")
	rsp := doToAppsResponse(0, do)
	rsp.AppGitName = appsRepository.ToGitEntity(do.AppGit).Name
	return rsp
}

// SyncWorkflows 同步工作流文件
// @post syncWorkflows
// @filter application.Jwt
func SyncWorkflows(appName string, appsRepository apps.Repository, gitDevice apps.IGitDevice) {
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	c := make(chan string, 100)
	gitEO := appsRepository.ToGitEntity(do.AppGit)
	if !gitDevice.PullWorkflows(do.GetWorkflowsRoot(), gitEO.Branch, gitEO.GetAuthHub(), c) {
		lstLog := collections.NewListFromChan(c)
		exception.ThrowWebExceptionf(403, "同步工作流文件失败:<br />%s", lstLog.ToString("<br />"))
	}
}

func doToAppsResponse(clusterId int64, do apps.DomainObject) response.AppsResponse {
	if clusterId == 0 {
		for i := range do.ClusterVer {
			clusterId = i
			break
		}
	}
	vo, exists := do.ClusterVer[clusterId]
	if !exists {
		vo = &apps.ClusterVerVO{}
	}
	return response.AppsResponse{
		AppName:           do.AppName,
		ActiveInstance:    do.ActiveInstance,
		DockerVer:         do.DockerVer,
		DockerImage:       do.DockerImage,
		ClusterVer:        *vo,
		AppGit:            do.AppGit,
		FrameworkGits:     do.FrameworkGits,
		DockerfilePath:    do.DockerfilePath,
		DockerNodeRole:    do.DockerNodeRole,
		DockerReplicas:    do.DockerReplicas,
		AdditionalScripts: do.AdditionalScripts,
		IsHealth:          len(do.ActiveInstance) >= do.DockerReplicas,
	}
}
