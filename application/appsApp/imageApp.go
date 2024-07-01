// @area /apps/
package appsApp

import (
	"context"
	"fops/domain/apps"
	"fops/domain/apps/event"
	"fops/domain/cluster"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/trace"
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

	// 如果仓库和集群的版本一致时，不允许同步
	if do.ClusterVer[clusterId] != nil && do.DockerVer == do.ClusterVer[clusterId].DockerVer {
		exception.ThrowWebExceptionfBool(clusterDO.IsNil(), 403, "版本一致，不需要同步")
	}

	// 先登陆仓库
	if dockerUserName != "" && dockerUserPwd != "" {
		container.Resolve[trace.IManager]().TraceHand("登陆镜像仓库").Run(func() {
			c := make(chan string, 100)
			if !appsIDockerDevice.Login(dockerHub, dockerUserName, dockerUserPwd, c) {
				lstLog := collections.NewListFromChan(c)
				exception.ThrowWebExceptionf(403, "镜像登陆失败:<br />%s", lstLog.ToString("<br />"))
			}
		})
	}

	c := make(chan string, 100)
	// 先拉取镜像
	container.Resolve[trace.IManager]().TraceHand("先拉取镜像").Run(func() {
		appsIDockerDevice.Pull(do.DockerImage, c)
	})

	// 首次创建还是更新镜像
	if appsIDockerSwarmDevice.ExistsDocker(appName) {
		// 更新镜像
		container.Resolve[trace.IManager]().TraceHand("更新镜像").Run(func() {
			if !appsIDockerSwarmDevice.SetImages(clusterDO, appName, do.DockerImage, c) {
				lstLog := collections.NewListFromChan(c)
				exception.ThrowWebExceptionf(403, "更新镜像失败:<br />%s", lstLog.ToString("<br />"))
			}
		})
	} else {
		// 创建容器服务
		container.Resolve[trace.IManager]().TraceHand("创建容器服务").Run(func() {
			if !appsIDockerSwarmDevice.CreateService(appName, do.DockerNodeRole, do.AdditionalScripts, clusterDO.DockerNetwork, do.DockerReplicas, do.DockerImage, do.LimitCpus, do.LimitMemory, c, context.Background()) {
				lstLog := collections.NewListFromChan(c)
				exception.ThrowWebExceptionf(403, "创建容器服务失败:<br />%s", lstLog.ToString("<br />"))
			}
		})
	}

	// 更新集群版本信息
	do.UpdateBuildVer(true, clusterId, 0)
	_, _ = appsRepository.UpdateClusterVer(appName, do.ClusterVer)
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
func RestartDocker(clusterId int64, appName string, appsIDockerSwarmDevice apps.IDockerSwarmDevice, clusterRepository cluster.Repository, appsRepository apps.Repository) {
	clusterDO := clusterRepository.ToEntity(clusterId)
	exception.ThrowWebExceptionfBool(clusterDO.IsNil(), 403, "集群不存在")

	c := make(chan string, 100)
	if !appsIDockerSwarmDevice.Restart(clusterDO, appName, c) {
		// 重启失败时，判断容器是否存在
		if !appsIDockerSwarmDevice.ExistsDocker(appName) {
			c = make(chan string, 100)
			// 创建容器服务
			do := appsRepository.ToEntity(appName)
			if !appsIDockerSwarmDevice.CreateService(appName, do.DockerNodeRole, do.AdditionalScripts, clusterDO.DockerNetwork, do.DockerReplicas, do.GetCurClusterDockerImage(clusterId), do.LimitCpus, do.LimitMemory, c, context.Background()) {
				lstLog := collections.NewListFromChan(c)
				exception.ThrowWebExceptionf(403, "创建容器服务失败:<br />%s", lstLog.ToString("<br />"))
			}
		} else {
			lstLog := collections.NewListFromChan(c)
			exception.ThrowWebExceptionf(403, "容器重启失败:<br />%s", lstLog.ToString("<br />"))
		}
	}
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
