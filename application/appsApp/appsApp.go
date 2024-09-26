// @area /apps/
package appsApp

import (
	"fmt"
	"fops/application/appsApp/request"
	"fops/application/appsApp/response"
	"fops/domain/apps"
	"fops/domain/cluster"
	"fops/domain/fSchedule"
	"fops/domain/logData"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/utils/file"
	"path/filepath"
	"strconv"
	"strings"
)

// Add 添加应用
// @post add
// @filter application.Jwt
func Add(req request.AddRequest, appsRepository apps.Repository) {
	do := mapper.Single[apps.DomainObject](req)
	do.IsSys = false
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
func Update(req request.UpdateRequest, appsRepository apps.Repository) {
	do := appsRepository.ToEntity(req.AppName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	client := docker.NewClient()
	if exists, err := client.Service.Exists(req.AppName); exists || err != nil {
		// 更新镜像
		if (req.ClusterDockerImage != "" && do.ClusterVer[req.ClusterId] != nil && req.ClusterDockerImage != do.ClusterVer[req.ClusterId].DockerImage) || do.DockerReplicas != req.DockerReplicas {
			err = client.Service.SetImagesAndReplicas(req.AppName, req.ClusterDockerImage, req.DockerReplicas)
			exception.ThrowWebExceptionError(403, err)
		} else if do.DockerReplicas != req.DockerReplicas {
			// 更新副本数量
			err = client.Service.SetReplicas(req.AppName, req.DockerReplicas)
			exception.ThrowWebExceptionError(403, err)
		}
	}

	// 更新应用信息
	newDO := mapper.Single[apps.DomainObject](req, func(newVal *apps.DomainObject) {
		newVal.ClusterVer = do.ClusterVer
	})

	// 删除末尾的/
	if strings.HasSuffix(newDO.AdditionalScripts, "\\") {
		newDO.AdditionalScripts = newDO.AdditionalScripts[:len(newDO.AdditionalScripts)-1]
	}

	// 更新部署的镜像
	if newDO.ClusterVer[req.ClusterId] != nil && req.ClusterDockerImage != "" {
		newDO.ClusterVer[req.ClusterId].DockerImage = req.ClusterDockerImage
		newDO.ClusterVer[req.ClusterId].DeploySuccessAt = dateTime.Now()
		if strings.Contains(req.ClusterDockerImage, ":") {
			newDO.ClusterVer[req.ClusterId].DockerVer = parse.ToInt(strings.Split(req.ClusterDockerImage, ":")[1])
		}
	}

	err := appsRepository.UpdateApp(newDO)
	exception.ThrowWebExceptionError(403, err)
}

// Delete 删除应用
// @post delete
// @filter application.Jwt
func Delete(appName string, appsRepository apps.Repository) {
	exception.ThrowWebExceptionBool(strings.Trim(appName, "") == "", 403, "参数不完整")
	// 删除服务
	client := docker.NewClient()
	exists, err := client.Service.Exists(appName)
	// 当err!=nil时，也认为服务是存在的。
	if exists || err != nil {
		err = client.Service.Delete(appName)
		exception.ThrowWebExceptionError(403, err)
	}

	// 删除应用
	_, err = appsRepository.Delete(appName)
	exception.ThrowWebExceptionError(403, err)
}

// DropDownList 应用列表
// @post dropDownList
// @filter application.Jwt
func DropDownList(isAll bool, appsRepository apps.Repository) collections.List[apps.ShortEO] {
	return appsRepository.ToShortList(isAll)
}

// List 应用列表
// @post list
// @filter application.Jwt
func List(clusterId int64, isSys bool, appsRepository apps.Repository, logDataRepository logData.Repository, clusterRepository cluster.Repository, fScheduleHttp fSchedule.Http) collections.List[response.AppsResponse] {
	lstGit := appsRepository.ToGitListAll(-1)
	countList := logDataRepository.StatCount()
	clusterDO := clusterRepository.ToEntity(clusterId)
	var taskGroupStatList collections.List[fSchedule.StatTaskEO]

	// 获取任务组的数据统计
	if clusterDO.FScheduleAddr != "" {
		taskGroupStatList = fScheduleHttp.StatList(clusterDO.FScheduleAddr)
	}

	lst := collections.NewList[response.AppsResponse]()
	appsRepository.ToListBySys(isSys).Foreach(func(item *apps.DomainObject) {
		appsResponse := doToAppsResponse(clusterId, *item)
		// Git名称
		appsResponse.AppGitName = lstGit.Where(func(gitItem apps.GitEO) bool {
			return item.AppGit == parse.ToInt64(gitItem.Id)
		}).First().Name

		// 统计日志数量
		countList.Foreach(func(logItem *logData.LogCountEO) {
			if item.AppName != logItem.AppName {
				return
			}
			switch logItem.LogLevel {
			case eumLogLevel.Error: // 日志异常数量
				appsResponse.LogErrorCount++
			case eumLogLevel.Warning: // 日志警告数量
				appsResponse.LogWaringCount++
			default:
			}
		})

		// 统计任务组执行数量
		taskGroupStatList.Foreach(func(statTaskEO *fSchedule.StatTaskEO) {
			if statTaskEO.ClientName != item.AppName {
				return
			}
			switch statTaskEO.ExecuteStatus {
			case 3: // 任务组执行失败数量
				appsResponse.TaskFailCount++
			case 2: // 任务组执行成功数量
				appsResponse.TaskSuccessCount++
			}
		})

		// 获取工作流文件名称
		workflowsNames := file.GetFiles(item.GetWorkflowsDir(), "*.yml", true)
		for _, workflowsYmlPath := range workflowsNames {
			// 读取工作流内容
			workflowsYml := configure.NewYamlConfig("")
			_ = workflowsYml.LoadContent([]byte(file.ReadString(workflowsYmlPath)))

			// 取出集群ID
			clusterIds, _ := workflowsYml.GetArray("jobs.clusterId")
			if len(clusterIds) == 0 {
				clusterIds = []any{clusterId}
			}

			// 只筛选出对应集群Id的工作流
			for _, cId := range clusterIds {
				if clusterId == parse.ToInt64(cId) {
					workflowsYmlPath = filepath.Base(workflowsYmlPath)
					appsResponse.WorkflowsNames = append(appsResponse.WorkflowsNames, workflowsYmlPath[:strings.Index(workflowsYmlPath, ".yml")])
				}
			}
		}

		// 容器资源占用统计
		appsResponse.CpuUsagePercent = item.DockerInspect.Where(func(item apps.DockerInspectVO) bool {
			return item.CpuUsagePercent > 0
		}).Average(func(item apps.DockerInspectVO) any {
			return item.CpuUsagePercent
		})
		appsResponse.CpuUsagePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", appsResponse.CpuUsagePercent), 64)

		appsResponse.MemoryUsagePercent = item.DockerInspect.Where(func(item apps.DockerInspectVO) bool {
			return item.MemoryUsagePercent > 0
		}).Average(func(item apps.DockerInspectVO) any {
			return item.MemoryUsagePercent
		})
		appsResponse.MemoryUsagePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", appsResponse.MemoryUsagePercent), 64)

		appsResponse.MemoryUsage = parse.ToUInt64(item.DockerInspect.Where(func(item apps.DockerInspectVO) bool {
			return item.MemoryUsage > 0
		}).Average(func(item apps.DockerInspectVO) any {
			return item.MemoryUsage
		}))
		appsResponse.MemoryLimit = parse.ToUInt64(item.DockerInspect.Where(func(item apps.DockerInspectVO) bool {
			return item.MemoryLimit > 0
		}).Average(func(item apps.DockerInspectVO) any {
			return item.MemoryLimit
		}))

		lst.Add(appsResponse)
	})
	return lst
}

// Info 查询应用
// @post info
// @filter application.Jwt
func Info(clusterId int64, appName string, appsRepository apps.Repository) response.AppsResponse {
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")
	rsp := doToAppsResponse(clusterId, do)
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
		DockerInstances:   do.DockerInstances,
		DockerVer:         do.DockerVer,
		DockerImage:       do.DockerImage,
		ClusterVer:        *vo,
		AppGit:            do.AppGit,
		FrameworkGits:     do.FrameworkGits,
		DockerfilePath:    do.DockerfilePath,
		DockerNodeRole:    do.DockerNodeRole,
		DockerReplicas:    do.DockerReplicas,
		AdditionalScripts: do.AdditionalScripts,
		IsHealth:          do.DockerInstances >= do.DockerReplicas,
		LimitCpus:         do.LimitCpus,
		LimitMemory:       do.LimitMemory,
	}
}
