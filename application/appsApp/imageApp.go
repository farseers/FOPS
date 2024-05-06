// @area /apps/
package appsApp

import (
	"context"
	"fops/domain/apps"
	"fops/domain/apps/event"
	"fops/domain/cluster"
	"github.com/farseer-go/collections"
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

	c := make(chan string, 100)
	// 先拉取镜像
	appsIDockerDevice.Pull(do.DockerImage, c)

	// 首次创建还是更新镜像
	if appsIDockerSwarmDevice.ExistsDocker(appName) {
		// 更新镜像
		if !appsIDockerSwarmDevice.SetImages(clusterDO, appName, do.DockerImage, do.DockerReplicas, c) {
			lstLog := collections.NewListFromChan(c)
			exception.ThrowWebExceptionf(403, "同步仓库版本失败:<br />%s", lstLog.ToString("<br />"))
		}
	} else {
		// 创建容器服务
		if !appsIDockerSwarmDevice.CreateService(appName, do.DockerNodeRole, do.AdditionalScripts, clusterDO.DockerNetwork, do.DockerReplicas, do.DockerImage, c, context.Background()) {
			lstLog := collections.NewListFromChan(c)
			exception.ThrowWebExceptionf(403, "同步仓库版本失败:<br />%s", lstLog.ToString("<br />"))
		}
	}

	// 更新集群版本信息
	do.UpdateBuildVer(true, clusterId, 0)
	_, _ = appsRepository.UpdateClusterVer(appName, do.ClusterVer)
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

// UpdateDockerImage 更新仓库版本
// @post updateDockerImage
func UpdateDockerImage(appName string, dockerImage string, buildNumber int, clusterId int64, dockerHub, dockerUserName, dockerUserPwd string, appsIDockerDevice apps.IDockerDevice, appsIDockerSwarmDevice apps.IDockerSwarmDevice, appsRepository apps.Repository, clusterRepository cluster.Repository) {
	// 更新仓库版本
	event.DockerPushedEvent{BuildNumber: buildNumber, AppName: appName, ImageName: dockerImage}.PublishEvent()

	// 如果集群ID大于0，则同步应用
	if clusterId > 0 {
		// 先登陆仓库
		if dockerUserName != "" && dockerUserPwd != "" {
			c := make(chan string, 100)
			if !appsIDockerDevice.Login(dockerHub, dockerUserName, dockerUserPwd, c) {
				lstLog := collections.NewListFromChan(c)
				exception.ThrowWebExceptionf(403, "镜像登陆失败:<br />%s", lstLog.ToString("<br />"))
			}
		}

		// 同步镜像
		SyncDockerImage(clusterId, appName, appsIDockerSwarmDevice, appsIDockerDevice, appsRepository, clusterRepository)
	}
}

// ClearDockerImage 清除Docker镜像
// @post build/clearDockerImage
// @filter application.Jwt
func ClearDockerImage(device apps.IDockerDevice) {
	c := make(chan string, 100)
	device.ClearImages(c)
}

// RestartDocker 重启容器
// @post build/restartDocker
// @filter application.Jwt
func RestartDocker(clusterId int64, appName string, appsIDockerSwarmDevice apps.IDockerSwarmDevice, clusterRepository cluster.Repository) {
	clusterDO := clusterRepository.ToEntity(clusterId)
	exception.ThrowWebExceptionfBool(clusterDO.IsNil(), 403, "集群不存在")

	c := make(chan string, 100)
	if !appsIDockerSwarmDevice.Restart(clusterDO, appName, c) {
		lstLog := collections.NewListFromChan(c)
		exception.ThrowWebExceptionf(403, "容器重启失败:<br />%s", lstLog.ToString("<br />"))
	}
}
