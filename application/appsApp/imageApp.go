// @area /apps/
package appsApp

import (
	"context"
	"fops/domain/_/eumBuildStatus"
	"fops/domain/apps"
	"fops/domain/apps/event"
	"fops/domain/cluster"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"strings"
)

// SyncDockerImage 同步仓库版本
// @post build/syncDockerImage
// @filter application.Jwt
func SyncDockerImage(clusterId int64, appName string, appsIDockerSwarmDevice apps.IDockerSwarmDevice, appsIDockerDevice apps.IDockerDevice, appsRepository apps.Repository, clusterRepository cluster.Repository) {
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	clusterDO := clusterRepository.ToEntity(clusterId)
	exception.ThrowWebExceptionfBool(clusterDO.IsNil(), 403, "集群不存在")

	// 如果仓库和集群的版本一致时，不允许同步
	if do.ClusterVer[clusterId] != nil && do.DockerVer == do.ClusterVer[clusterId].DockerVer {
		exception.ThrowWebExceptionfBool(clusterDO.IsNil(), 403, "版本一致，不需要同步")
	}

	// 先登陆仓库
	if clusterDO.DockerUserName != "" && clusterDO.DockerUserPwd != "" {
		c := make(chan string, 100)
		if !appsIDockerDevice.Login(clusterDO.DockerHub, clusterDO.DockerUserName, clusterDO.DockerUserPwd, c) {
			lstLog := collections.NewListFromChan(c)
			exception.ThrowWebExceptionf(403, "镜像登陆失败:<br />%s", lstLog.ToString("<br />"))
		}
	}

	c := make(chan string, 100)
	// 先拉取镜像
	appsIDockerDevice.Pull(do.DockerImage, c)

	// 首次创建还是更新镜像
	if appsIDockerSwarmDevice.ExistsDocker(appName) {
		// 更新镜像
		if !appsIDockerSwarmDevice.SetImages(clusterDO, appName, do.DockerImage, c) {
			lstLog := collections.NewListFromChan(c)
			exception.ThrowWebExceptionf(403, "同步仓库版本失败:<br />%s", lstLog.ToString("<br />"))
		}
	} else {
		// 创建容器服务
		if !appsIDockerSwarmDevice.CreateService(appName, do.DockerNodeRole, do.AdditionalScripts, clusterDO.DockerNetwork, do.DockerReplicas, do.DockerImage, do.LimitCpus, do.LimitMemory, c, context.Background()) {
			lstLog := collections.NewListFromChan(c)
			exception.ThrowWebExceptionf(403, "同步仓库版本失败:<br />%s", lstLog.ToString("<br />"))
		}
	}

	// 更新集群版本信息
	do.UpdateBuildVer(true, clusterId, 0)
	_, _ = appsRepository.UpdateClusterVer(appName, do.ClusterVer)

	// 更新构建中状态的构建记录
	_, _ = appsRepository.UpdateFailDockerImage(appName, do.DockerImage)
}

// UpdateDockerImage 更新仓库版本
// @post updateDockerImage
func UpdateDockerImage(clusterId int64, appName string, dockerImage string, buildNumber int, dockerHub, dockerUserName, dockerUserPwd string, appsIDockerDevice apps.IDockerDevice, appsIDockerSwarmDevice apps.IDockerSwarmDevice, appsRepository apps.Repository, clusterRepository cluster.Repository) {
	buildLogEO := apps.BuildEO{
		ClusterId:     clusterId,
		BuildNumber:   buildNumber,
		Status:        eumBuildStatus.Finish,
		CreateAt:      dateTime.Now(),
		BuildServerId: core.AppId,
		AppName:       appName,
		WorkflowsName: "远程",
		DockerImage:   dockerImage,
	}

	defer func() {
		// 手动创建一个构建记录
		buildLogEO.FinishAt = dateTime.Now()
		_ = appsRepository.AddBuild(buildLogEO)
	}()

	// 更新仓库版本
	event.DockerPushedEvent{BuildNumber: buildNumber, AppName: appName, ImageName: dockerImage}.PublishEvent()

	// 如果集群ID大于0，则同步应用
	if clusterId < 1 {
		return
	}

	// 同步镜像
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	clusterDO := clusterRepository.ToEntity(clusterId)
	exception.ThrowWebExceptionfBool(clusterDO.IsNil(), 403, "集群不存在")

	// 先登陆仓库
	if dockerUserName != "" && dockerUserPwd != "" {
		c := make(chan string, 100)
		if !appsIDockerDevice.Login(dockerHub, dockerUserName, dockerUserPwd, c) {
			lstLog := collections.NewListFromChan(c)
			exception.ThrowWebExceptionf(403, "镜像登陆失败:\r\n%s", lstLog.ToString("\n"))
		}
	}

	c := make(chan string, 100)
	// 先拉取镜像
	appsIDockerDevice.Pull(do.DockerImage, c)

	// 先判断容器服务是否存在（运行中）
	if appsIDockerSwarmDevice.ExistsDocker(appName) {
		// 更新镜像
		if !appsIDockerSwarmDevice.SetImages(clusterDO, appName, do.DockerImage, c) {
			lstLog := collections.NewListFromChan(c)
			exception.ThrowWebExceptionf(403, "更新镜像失败:\n%s", lstLog.ToString("\n"))
		}
	} else {
		// 创建容器服务
		if !appsIDockerSwarmDevice.CreateService(appName, do.DockerNodeRole, do.AdditionalScripts, clusterDO.DockerNetwork, do.DockerReplicas, do.DockerImage, do.LimitCpus, do.LimitMemory, c, context.Background()) {
			lstLog := collections.NewListFromChan(c)
			exception.ThrowWebExceptionf(403, "创建容器服务失败:<br />%s", lstLog.ToString("<br />"))
		}
	}

	// 更新集群版本信息
	do.UpdateBuildVer(true, clusterId, 0)
	_, _ = appsRepository.UpdateClusterVer(appName, do.ClusterVer)

	buildLogEO.IsSuccess = true
}

// ClearDockerImage 清除Docker镜像
// @post build/clearDockerImage
// @filter application.Jwt
func ClearDockerImage(dockerDevice apps.IDockerDevice) {
	dockerDevice.Kill("FOPS-Build")
	c := make(chan string, 100)
	dockerDevice.ClearImages(c)
}

// RestartDocker 重启容器
// @post restartDocker
// @filter application.Jwt
func RestartDocker(clusterId int64, appName string, appsIDockerSwarmDevice apps.IDockerSwarmDevice, clusterRepository cluster.Repository, appsRepository apps.Repository) {
	c := make(chan string, 100)
	if !appsIDockerSwarmDevice.Restart(appName, c) {
		// 重启失败时，判断容器是否存在
		if !appsIDockerSwarmDevice.ExistsDocker(appName) {
			clusterDO := clusterRepository.ToEntity(clusterId)
			if clusterDO.IsNil() {
				clusterDO.DockerNetwork = "net"
			}

			c = make(chan string, 100)
			// 创建容器服务
			do := appsRepository.ToEntity(appName)
			dockerImage := do.GetCurClusterDockerImage(clusterId)
			if dockerImage == "" {
				exception.ThrowWebExceptionf(403, "该集群没有可用的镜像")
			}
			if !appsIDockerSwarmDevice.CreateService(appName, do.DockerNodeRole, do.AdditionalScripts, clusterDO.DockerNetwork, do.DockerReplicas, dockerImage, do.LimitCpus, do.LimitMemory, c, context.Background()) {
				lstLog := collections.NewListFromChan(c)
				exception.ThrowWebExceptionf(403, "创建容器服务失败:<br />%s", lstLog.ToString("<br />"))
			}
		} else {
			lstLog := collections.NewListFromChan(c)
			exception.ThrowWebExceptionf(403, "容器重启失败:<br />%s", lstLog.ToString("<br />"))
		}
	}
}

// SetReplicas 更新副本实例数量
// @post setReplicas
// @filter application.Jwt
func SetReplicas(appName string, dockerReplicas int, appsRepository apps.Repository, appsIDockerSwarmDevice apps.IDockerSwarmDevice) {
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	// 更新副本数量
	c := make(chan string, 100)
	if !appsIDockerSwarmDevice.SetReplicas(cluster.DomainObject{}, appName, dockerReplicas, c) {
		lstLog := collections.NewListFromChan(c)
		exception.ThrowWebExceptionf(403, "更新副本失败:<br />%s", lstLog.ToString("<br />"))
	}

	do.DockerReplicas = dockerReplicas
	err := appsRepository.UpdateApp(do)
	exception.ThrowWebExceptionError(403, err)
}

// DeleteService 删除容器服务
// @post deleteService
// @filter application.Jwt
func DeleteService(appName string, appsRepository apps.Repository, appsIDockerSwarmDevice apps.IDockerSwarmDevice) {
	exception.ThrowWebExceptionBool(strings.Trim(appName, "") == "", 403, "参数不完整")
	// 删除服务
	c := make(chan string, 100)
	if !appsIDockerSwarmDevice.DeleteService(appName, c) {
		lstLog := collections.NewListFromChan(c)
		exception.ThrowWebExceptionf(403, "删除容器服务失败:<br />%s", lstLog.ToString("<br />"))
	}
}
