// @area /apps/
package appsApp

import (
	"fops/domain/_/eumBuildStatus"
	"fops/domain/apps"
	"fops/domain/cluster"
	"strings"

	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
)

// UpdateDockerImage 更新仓库版本
// @post updateDockerImage
func UpdateDockerImage(appName string, dockerImage string, updateDelay int, buildNumber int, dockerHub, dockerUserName, dockerUserPwd string, appsRepository apps.Repository, clusterRepository cluster.Repository) {
	clusterDO := clusterRepository.GetLocalCluster()
	exception.ThrowWebExceptionfBool(clusterDO.IsNil(), 403, "集群不存在")

	buildLogEO := apps.BuildEO{
		ClusterId:     clusterDO.Id,
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
		_ = appsRepository.AddBuild(&buildLogEO)
	}()

	// 同步镜像
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	// 先登陆仓库
	wait := docker.DefaultClient.Hub.Login(dockerHub, dockerUserName, dockerUserPwd)
	result, exitCode := wait.WaitToList()
	exception.ThrowRefuseExceptionBool(exitCode != 0, "镜像登陆失败:"+result.ToString(","))

	// 先拉取镜像
	wait = docker.DefaultClient.Images.Pull(dockerImage)
	result, exitCode = wait.WaitToList()
	exception.ThrowRefuseExceptionBool(exitCode != 0, result.ToString(","))

	// 服务存在，才更新，否则自动创建
	if !createService(docker.DefaultClient, appName, dockerImage, appsRepository, clusterRepository) {
		// 检查并更新配置版本（这里不需要返回值，因为后续的 SetImages 会触发服务更新）
		_ = docker.DefaultClient.SyncConfig(appName, "/app/config.yaml")

		// 更新镜像
		wait = docker.DefaultClient.Service.SetImages(appName, dockerImage, updateDelay)
		result, exitCode = wait.WaitToList()
		exception.ThrowRefuseExceptionBool(exitCode != 0, result.ToString(","))
	}

	// 更新集群版本信息
	do.UpdateBuildVer(true, clusterDO.Id, 0, buildNumber, dockerImage)
	_, _ = appsRepository.UpdateClusterVer(appName, do.ClusterVer)

	buildLogEO.IsSuccess = true
}

// ClearDockerImage 清除Docker镜像
// @post build/clearDockerImage
// @filter application.Jwt
func ClearDockerImage() {
	docker.DefaultClient.Container.Kill("FOPS-Build")
	docker.DefaultClient.Container.Kill("FOPS-AutoBuild")
	docker.DefaultClient.Images.ClearImages()
}

// RestartDocker 重启容器
// @post restartDocker
// @filter application.Jwt
func RestartDocker(appName string, clusterRepository cluster.Repository, appsRepository apps.Repository) {
	// 服务存在，才重启，否则自动创建
	if !createService(docker.DefaultClient, appName, "", appsRepository, clusterRepository) {
		// 检查并更新配置版本（更新配置会自动触发服务滚动更新）
		configUpdated := docker.DefaultClient.SyncConfig(appName, "/app/config.yaml")

		// 如果配置没有更新，才执行重启
		if !configUpdated {
			wait := docker.DefaultClient.Service.Restart(appName)
			result, exitCode := wait.WaitToList()
			exception.ThrowRefuseExceptionBool(exitCode != 0, result.ToString(","))
		}
	}
}

// SetReplicas 更新副本实例数量
// @post setReplicas
// @filter application.Jwt
func SetReplicas(appName string, dockerReplicas int, appsRepository apps.Repository) {
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	exists := docker.DefaultClient.Service.Exists(appName)

	// 更新副本数量
	if exists {
		wait := docker.DefaultClient.Service.SetReplicas(appName, dockerReplicas)
		result, exitCode := wait.WaitToList()
		exception.ThrowRefuseExceptionBool(exitCode != 0, result.ToString(","))
	}

	do.DockerReplicas = dockerReplicas
	err := appsRepository.UpdateApp(do)
	exception.ThrowWebExceptionError(403, err)
}

// DeleteService 删除容器服务
// @post deleteService
// @filter application.Jwt
func DeleteService(appName string, appsRepository apps.Repository) {
	exception.ThrowWebExceptionBool(strings.Trim(appName, "") == "", 403, "参数不完整")
	// 删除服务
	err := docker.DefaultClient.Service.Delete(appName)
	if err != nil {
		exception.ThrowWebExceptionError(403, err)
	}

	// 验证
	exists := docker.DefaultClient.Service.Exists(appName)
	exception.ThrowWebExceptionBool(exists, 403, "服务删除失败")
}

func createService(client *docker.Client, appName, dockerImage string, appsRepository apps.Repository, clusterRepository cluster.Repository) bool {
	clusterDO := clusterRepository.GetLocalCluster()
	// 服务不存在，则创建
	exists := client.Service.Exists(appName)
	if !exists {
		// 创建容器服务
		do := appsRepository.ToEntity(appName)
		if dockerImage == "" {
			dockerImage = do.GetCurClusterDockerImage(clusterDO.Id)
		}
		if dockerImage == "" {
			exception.ThrowWebExceptionf(403, "该集群没有可用的镜像")
		}

		if clusterDO.IsNil() {
			clusterDO.DockerNetwork = "net"
		}

		// 准备配置文件
		configVersion, err := client.Config.GetLastVersion(appName)
		exception.ThrowRefuseExceptionError(err)

		wait := client.Service.Create(appName, do.DockerNodeRole, do.AdditionalScripts, clusterDO.DockerNetwork, do.DockerReplicas, dockerImage, do.LimitCpus, do.LimitMemory, docker.ConfigTarget{
			Name:   configVersion.Spec.Name,
			Target: "/app/config.yaml",
		})
		result, exitCode := wait.WaitToList()
		exception.ThrowRefuseExceptionBool(exitCode != 0, result.ToString(","))
		return true
	}
	return false
}
