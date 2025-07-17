// @area /apps/
package appsApp

import (
	"context"
	"fmt"
	"fops/application/appsApp/request"
	"fops/application/appsApp/response"
	"fops/domain"
	"fops/domain/apps"
	"fops/domain/cluster"
	"fops/domain/clusterNode"
	"fops/domain/fSchedule"
	"fops/domain/logData"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/core/eumLogLevel"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/utils/file"
)

// @summary 添加应用
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
func Update(req request.UpdateRequest, appsRepository apps.Repository, clusterRepository cluster.Repository) {
	do := appsRepository.ToEntity(req.AppName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")
	clusterDO := clusterRepository.GetLocalCluster()

	client := docker.NewClient()
	if exists, err := client.Service.Exists(req.AppName); exists || err != nil {
		// 更新镜像
		if (req.ClusterDockerImage != "" && req.ClusterDockerImage != do.ClusterVer.GetValue(clusterDO.Id).DockerImage) || do.DockerReplicas != req.DockerReplicas {
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
	newDO.AdditionalScripts = strings.TrimSuffix(newDO.AdditionalScripts, "\\")

	// 更新部署的镜像
	if newDO.ClusterVer.ContainsKey(clusterDO.Id) && req.ClusterDockerImage != "" {
		clusterVerVO := newDO.ClusterVer.GetValue(clusterDO.Id)
		clusterVerVO.DockerImage = req.ClusterDockerImage
		clusterVerVO.DeploySuccessAt = dateTime.Now()
		if strings.Contains(req.ClusterDockerImage, ":") {
			clusterVerVO.DockerVer = parse.ToInt(strings.Split(req.ClusterDockerImage, ":")[1])
		}
		newDO.ClusterVer.Add(clusterDO.Id, clusterVerVO)
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
	resList := appsRepository.ToShortList(isAll)
	// cluster_node 节点信息
	clusterNode.NodeList.Foreach(func(node *docker.DockerNodeVO) {
		resList.Add(apps.ShortEO{AppName: fmt.Sprintf("%s(%s)", node.IP, node.NodeName)})
	})
	return resList
}

// List 应用列表
// @post list
// @filter application.Jwt
func List(isSys bool, appsRepository apps.Repository, logDataRepository logData.Repository, clusterRepository cluster.Repository, fScheduleHttp fSchedule.Http) collections.List[response.AppsResponse] {
	countList := logDataRepository.StatCount()
	// 获取任务组的数据统计
	var taskGroupStatList collections.List[fSchedule.StatTaskEO]
	if clusterDO := clusterRepository.GetLocalCluster(); clusterDO.FScheduleAddr != "" {
		taskGroupStatList = fScheduleHttp.StatList(clusterDO.FScheduleAddr)
	}
	lstCluster := clusterRepository.ToList()
	lst := collections.NewList[response.AppsResponse]()
	appsRepository.ToListBySys(isSys).Foreach(func(item *apps.DomainObject) {
		// 在监控中心，副本数量=0的，不显示
		if isSys && item.DockerReplicas == 0 {
			return
		}
		appsResponse := doToAppsResponse(lstCluster, *item)
		// 统计日志数量
		countList.Foreach(func(logItem *logData.LogCountEO) {
			if item.AppName != logItem.AppName {
				return
			}
			switch logItem.LogLevel {
			case eumLogLevel.Error: // 日志异常数量
				appsResponse.LogErrorCount += logItem.LogCount
			case eumLogLevel.Warning: // 日志警告数量
				appsResponse.LogWaringCount += logItem.LogCount
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
				appsResponse.TaskFailCount += statTaskEO.Count
			case 2: // 任务组执行成功数量
				appsResponse.TaskSuccessCount += statTaskEO.Count
			}
		})

		// 获取工作流文件名称
		workflowsNames := file.GetFiles(item.GetWorkflowsDir(), "*.yml", true)
		for _, workflowsYmlPath := range workflowsNames {
			// 读取工作流内容
			workflowsYml := configure.NewYamlConfig("")
			_ = workflowsYml.LoadContent([]byte(file.ReadString(workflowsYmlPath)))

			// 取出限制的名称
			fopsName, _ := workflowsYml.Get("fopsName")
			// 只筛选出对应名称的工作流
			if fopsNameString := parse.ToString(fopsName); fopsNameString == "" || fopsNameString == item.AppName {
				// 判断账号权限
				strClusterId, _ := workflowsYml.Get("jobs.clusterId")
				clusterId := parse.ToInt(strClusterId)
				if domain.GetLoginAccount().ClusterIds.Contains(clusterId) {
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
func Info(appName string, appsRepository apps.Repository, clusterRepository cluster.Repository) response.AppsResponse {
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	localCluster := clusterRepository.GetLocalCluster()
	clusterVerVO := do.ClusterVer.GetValue(localCluster.Id)

	lstCluster := clusterRepository.ToList()
	appsResponse := doToAppsResponse(lstCluster, do)
	appsResponse.LocalClusterVer = response.ClusterVerVO{
		ClusterId:       clusterVerVO.ClusterId,
		ClusterName:     localCluster.Name,
		DockerImage:     clusterVerVO.DockerImage,
		DeploySuccessAt: clusterVerVO.DeploySuccessAt,
	}
	return appsResponse
}

// SyncWorkflows 同步工作流文件
// @post syncWorkflows
// @filter application.Jwt
func SyncWorkflows(appName string, appsRepository apps.Repository, gitDevice apps.IGitDevice) {
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	c := make(chan string, 100)
	gitEO := appsRepository.ToGitEntity(do.AppGit)
	if !gitDevice.PullWorkflows(context.Background(), do.GetWorkflowsRoot(), gitEO.Branch, gitEO.GetAuthHub(), c) {
		lstLog := collections.NewListFromChan(c)
		exception.ThrowWebExceptionf(403, "同步工作流文件失败:<br />%s", lstLog.ToString("<br />"))
	}
}

func doToAppsResponse(lstCluster collections.List[cluster.DomainObject], do apps.DomainObject) response.AppsResponse {
	clusterVer := collections.NewList[response.ClusterVerVO]()
	do.ClusterVer.Values().Foreach(func(clusterVerVO *apps.ClusterVerVO) {
		if curCluster := lstCluster.Find(func(item *cluster.DomainObject) bool {
			return item.Id == clusterVerVO.ClusterId
		}); curCluster != nil {
			vo := mapper.Single[response.ClusterVerVO](clusterVerVO)
			vo.ClusterName = curCluster.Name
			clusterVer.Add(vo)
		}
	})
	clusterVer = clusterVer.OrderBy(func(item response.ClusterVerVO) any {
		return item.ClusterId
	}).ToList()
	return response.AppsResponse{
		AppName:           do.AppName,
		AppGit:            do.AppGit,
		DockerInstances:   do.DockerInstances,
		ClusterVer:        clusterVer,
		FrameworkGits:     do.FrameworkGits,
		DockerNodeRole:    do.DockerNodeRole,
		DockerReplicas:    do.DockerReplicas,
		IsHealth:          do.DockerInstances >= do.DockerReplicas,
		LimitCpus:         do.LimitCpus,
		LimitMemory:       do.LimitMemory,
		AdditionalScripts: do.AdditionalScripts,
		DockerfilePath:    do.DockerfilePath,
		UTWorkflowsName:   do.UTWorkflowsName,
	}
}
