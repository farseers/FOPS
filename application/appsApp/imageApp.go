// @area /apps/
package appsApp

import (
	"fmt"
	"fops/domain/_/eumBuildStatus"
	"fops/domain/apps"
	"fops/domain/cluster"
	"os"
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

	// 更新仓库版本
	//event.DockerPushedEvent{BuildNumber: buildNumber, AppName: appName, ImageName: dockerImage}.PublishEvent()

	defer func() {
		// 手动创建一个构建记录
		buildLogEO.FinishAt = dateTime.Now()
		_ = appsRepository.AddBuild(&buildLogEO)
	}()

	// 同步镜像
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	client := docker.NewClient()
	// 先登陆仓库
	wait := client.Hub.Login(dockerHub, dockerUserName, dockerUserPwd)
	result, exitCode := wait.WaitToList()
	exception.ThrowRefuseExceptionBool(exitCode != 0, "镜像登陆失败:"+result.ToString(","))

	// 先拉取镜像
	wait = client.Images.Pull(dockerImage)
	result, exitCode = wait.WaitToList()
	exception.ThrowRefuseExceptionBool(exitCode != 0, result.ToString(","))

	// 服务存在，才更新，否则自动创建
	if !createService(client, appName, dockerImage, appsRepository, clusterRepository) {
		// 检查并更新配置版本
		updateServiceConfigIfNeeded(client, appName, appsRepository)

		// 更新镜像
		wait = client.Service.SetImages(appName, dockerImage, updateDelay)
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
	client := docker.NewClient()
	client.Container.Kill("FOPS-Build")
	client.Container.Kill("FOPS-AutoBuild")
	client.Images.ClearImages()
}

// RestartDocker 重启容器
// @post restartDocker
// @filter application.Jwt
func RestartDocker(appName string, clusterRepository cluster.Repository, appsRepository apps.Repository) {
	client := docker.NewClient()
	// 服务存在，才重启，否则自动创建
	if !createService(client, appName, "", appsRepository, clusterRepository) {
		// 检查并更新配置版本
		updateServiceConfigIfNeeded(client, appName, appsRepository)

		wait := client.Service.Restart(appName)
		result, exitCode := wait.WaitToList()
		exception.ThrowRefuseExceptionBool(exitCode != 0, result.ToString(","))
	}
}

// SetReplicas 更新副本实例数量
// @post setReplicas
// @filter application.Jwt
func SetReplicas(appName string, dockerReplicas int, appsRepository apps.Repository) {
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	client := docker.NewClient()
	exists := client.Service.Exists(appName)

	// 更新副本数量
	if exists {
		wait := client.Service.SetReplicas(appName, dockerReplicas)
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
	client := docker.NewClient()
	err := client.Service.Delete(appName)
	if err != nil {
		exception.ThrowWebExceptionError(403, err)
	}

	// 验证
	exists := client.Service.Exists(appName)
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
		configName := ensureConfigExists(client, appName, do.ConfigVer, appsRepository)

		wait := client.Service.Create(appName, do.DockerNodeRole, do.AdditionalScripts, clusterDO.DockerNetwork, do.DockerReplicas, dockerImage, do.LimitCpus, do.LimitMemory, configName)
		result, exitCode := wait.WaitToList()
		exception.ThrowRefuseExceptionBool(exitCode != 0, result.ToString(","))
		return true
	}
	return false
}

// ensureConfigExists 确保配置文件存在，如果不存在则创建默认配置
func ensureConfigExists(client *docker.Client, appName string, configVer int, appsRepository apps.Repository) string {
	// 如果版本号为0，说明还没有配置，创建默认配置
	if configVer == 0 {
		// 读取默认配置模板
		defaultConfig, err := os.ReadFile("/home/code/FOPS/config.yaml")
		if err != nil {
			// 如果读取失败，使用空配置
			defaultConfig = []byte("")
		}

		// 创建配置
		configName := fmt.Sprintf("%s_config_v1", appName)
		labels := map[string]string{
			"owner_service": appName,
			"version":       "1",
		}

		_, err = client.Config.Create(configName, defaultConfig, labels)
		if err != nil {
			// 如果创建失败，返回空字符串（不挂载配置）
			return ""
		}

		// 更新应用的配置版本号
		do := appsRepository.ToEntity(appName)
		do.ConfigVer = 1
		_ = appsRepository.UpdateApp(do)

		return configName
	}

	// 如果已有配置版本，返回配置名称
	return fmt.Sprintf("%s_config_v%d", appName, configVer)
}

// updateServiceConfigIfNeeded 检查并更新服务的配置版本
func updateServiceConfigIfNeeded(client *docker.Client, appName string, appsRepository apps.Repository) {
	// 获取应用信息
	do := appsRepository.ToEntity(appName)
	if do.IsNil() || do.ConfigVer == 0 {
		return
	}

	// 获取服务当前使用的配置
	serviceInfo, err := client.Service.Inspect(appName)
	if err != nil {
		return
	}

	// 检查服务是否有配置
	if len(serviceInfo.Spec.TaskTemplate.ContainerSpec.Configs) == 0 {
		// 服务没有配置，需要添加
		configName := fmt.Sprintf("%s_config_v%d", appName, do.ConfigVer)
		_ = client.Service.UpdateServiceConfig(appName, configName, "/app/config.yaml")
		return
	}

	// 获取当前配置的版本
	currentConfig := serviceInfo.Spec.TaskTemplate.ContainerSpec.Configs[0]
	currentConfigInfo, err := client.Config.Inspect(currentConfig.ConfigID)
	if err != nil {
		return
	}

	// 比较版本号
	currentVersion := 0
	if v, ok := currentConfigInfo.Spec.Labels["version"]; ok {
		fmt.Sscanf(v, "%d", &currentVersion)
	}

	// 如果版本不一致，更新配置
	if currentVersion != do.ConfigVer {
		configName := fmt.Sprintf("%s_config_v%d", appName, do.ConfigVer)
		_ = client.Service.UpdateServiceConfig(appName, configName, "/app/config.yaml")
	}
}
