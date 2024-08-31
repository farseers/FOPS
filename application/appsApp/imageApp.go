// @area /apps/
package appsApp

import (
	"fops/domain/_/eumBuildStatus"
	"fops/domain/apps"
	"fops/domain/apps/event"
	"fops/domain/cluster"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"strings"
)

// SyncDockerImage 同步仓库版本
// @post build/syncDockerImage
// @filter application.Jwt
func SyncDockerImage(clusterId int64, appName string, appsRepository apps.Repository, clusterRepository cluster.Repository) {
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	clusterDO := clusterRepository.ToEntity(clusterId)
	exception.ThrowWebExceptionfBool(clusterDO.IsNil(), 403, "集群不存在")

	// 如果仓库和集群的版本一致时，不允许同步
	if do.ClusterVer[clusterId] != nil && do.DockerVer == do.ClusterVer[clusterId].DockerVer {
		exception.ThrowWebExceptionfBool(clusterDO.IsNil(), 403, "版本一致，不需要同步")
	}

	client := docker.NewClient()
	// 先登陆仓库
	err := client.Hub.Login(clusterDO.DockerHub, clusterDO.DockerUserName, clusterDO.DockerUserPwd)
	if err != nil {
		exception.ThrowWebExceptionf(403, "镜像登陆失败:%s", err.Error())
	}

	// 先拉取镜像
	err = client.Images.Pull(do.DockerImage)
	exception.ThrowRefuseExceptionError(err)

	// 服务存在，才更新，否则自动创建
	if !createService(client, clusterId, appName, do.DockerImage, appsRepository, clusterRepository) {
		// 更新镜像
		err = client.Service.SetImages(appName, do.DockerImage)
		exception.ThrowRefuseExceptionError(err)
	}

	// 更新集群版本信息
	do.UpdateBuildVer(true, clusterId, 0)
	_, _ = appsRepository.UpdateClusterVer(appName, do.ClusterVer)

	// 更新构建中状态的构建记录
	_, _ = appsRepository.UpdateFailDockerImage(appName, do.DockerImage)
}

// UpdateDockerImage 更新仓库版本
// @post updateDockerImage
func UpdateDockerImage(clusterId int64, appName string, dockerImage string, buildNumber int, dockerHub, dockerUserName, dockerUserPwd string, appsRepository apps.Repository, clusterRepository cluster.Repository) {
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

	// 更新仓库版本
	event.DockerPushedEvent{BuildNumber: buildNumber, AppName: appName, ImageName: dockerImage}.PublishEvent()

	// 如果集群ID大于0，则同步应用
	if clusterId < 1 {
		return
	}

	defer func() {
		// 手动创建一个构建记录
		buildLogEO.FinishAt = dateTime.Now()
		_ = appsRepository.AddBuild(buildLogEO)
	}()

	// 同步镜像
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	clusterDO := clusterRepository.ToEntity(clusterId)
	exception.ThrowWebExceptionfBool(clusterDO.IsNil(), 403, "集群不存在")

	client := docker.NewClient()

	// 先登陆仓库
	err := client.Hub.Login(dockerHub, dockerUserName, dockerUserPwd)
	if err != nil {
		exception.ThrowWebExceptionf(403, "镜像登陆失败:%s", err.Error())
	}

	// 先拉取镜像
	err = client.Images.Pull(do.DockerImage)
	exception.ThrowRefuseExceptionError(err)

	// 服务存在，才更新，否则自动创建
	if !createService(client, clusterId, appName, do.DockerImage, appsRepository, clusterRepository) {
		// 更新镜像
		err = client.Service.SetImages(appName, do.DockerImage)
		exception.ThrowRefuseExceptionError(err)
	}

	// 更新集群版本信息
	do.UpdateBuildVer(true, clusterId, 0)
	_, _ = appsRepository.UpdateClusterVer(appName, do.ClusterVer)

	buildLogEO.IsSuccess = true
}

// ClearDockerImage 清除Docker镜像
// @post build/clearDockerImage
// @filter application.Jwt
func ClearDockerImage() {
	client := docker.NewClient()
	client.Container.Kill("FOPS-Build")
	_ = client.Images.ClearImages()
}

// RestartDocker 重启容器
// @post restartDocker
// @filter application.Jwt
func RestartDocker(clusterId int64, appName string, clusterRepository cluster.Repository, appsRepository apps.Repository) {
	client := docker.NewClient()
	// 服务存在，才重启，否则自动创建
	if !createService(client, clusterId, appName, "", appsRepository, clusterRepository) {
		err := client.Service.Restart(appName)
		exception.ThrowRefuseExceptionError(err)
	}
}

// SetReplicas 更新副本实例数量
// @post setReplicas
// @filter application.Jwt
func SetReplicas(appName string, dockerReplicas int, appsRepository apps.Repository) {
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	client := docker.NewClient()
	exists, err := client.Service.Exists(appName)

	// 更新副本数量
	if exists || err != nil {
		err = client.Service.SetReplicas(appName, dockerReplicas)
		exception.ThrowWebExceptionError(403, err)
	}

	do.DockerReplicas = dockerReplicas
	err = appsRepository.UpdateApp(do)
	exception.ThrowWebExceptionError(403, err)
}

// DeleteService 删除容器服务
// @post deleteService
// @filter application.Jwt
func DeleteService(appName string, appsRepository apps.Repository) {
	exception.ThrowWebExceptionBool(strings.Trim(appName, "") == "", 403, "参数不完整")
	// 删除服务
	client := docker.NewClient()
	_ = client.Service.Delete(appName)
	// 验证
	exists, _ := client.Service.Exists(appName)
	exception.ThrowWebExceptionBool(exists, 403, "服务删除失败")
}

func createService(client *docker.Client, clusterId int64, appName, dockerImage string, appsRepository apps.Repository, clusterRepository cluster.Repository) bool {
	// 服务不存在，则创建
	exists, err := client.Service.Exists(appName)
	if !exists && err == nil {
		// 创建容器服务
		do := appsRepository.ToEntity(appName)
		if dockerImage == "" {
			dockerImage = do.GetCurClusterDockerImage(clusterId)
		}
		if dockerImage == "" {
			exception.ThrowWebExceptionf(403, "该集群没有可用的镜像")
		}

		clusterDO := clusterRepository.ToEntity(clusterId)
		if clusterDO.IsNil() {
			clusterDO.DockerNetwork = "net"
		}

		err = client.Service.Create(appName, do.DockerNodeRole, do.AdditionalScripts, clusterDO.DockerNetwork, do.DockerReplicas, dockerImage, do.LimitCpus, do.LimitMemory)
		exception.ThrowRefuseExceptionError(err)
		return true
	}
	return false
}
