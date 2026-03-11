// @area /apps/config/
package appsApp

import (
	"fmt"
	"fops/domain/apps"
	"os"

	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/exception"
)

// ConfigResponse 配置响应
type ConfigResponse struct {
	Content         string // 配置内容
	AppConfigVer    int    // 应用数据库中的配置版本号
	DockerConfigVer string // Docker中实际使用的配置版本号
}

// GetConfig 获取应用配置文件
// @get config/get
// @filter application.Jwt
func GetConfig(appName string, appsRepository apps.Repository) ConfigResponse {
	exception.ThrowWebExceptionBool(appName == "", 403, "应用名称不能为空")

	// 获取应用信息
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	client := docker.NewClient()

	response := ConfigResponse{
		AppConfigVer: do.ConfigVer,
	}

	// 尝试从 Docker Config 获取配置
	configInfo, err := client.Config.InspectByService(appName)
	if err == nil && configInfo.ID != "" {
		// 返回配置内容
		response.Content = configInfo.Spec.Data
		// 从 Labels 中获取 Docker Config 的版本号
		if version, ok := configInfo.Spec.Labels["version"]; ok {
			response.DockerConfigVer = version
		}
		return response
	}

	// 如果不存在，返回默认模板
	defaultConfig, err := os.ReadFile("./tpl.yaml")
	if err != nil {
		exception.ThrowWebExceptionf(403, "读取默认配置模板失败: %v", err)
	}

	response.Content = string(defaultConfig)
	response.DockerConfigVer = "未创建"
	return response
}

// SaveConfig 保存应用配置文件
// @post config/save
// @filter application.Jwt
func SaveConfig(appName string, content string, appsRepository apps.Repository) {
	exception.ThrowWebExceptionBool(appName == "", 403, "应用名称不能为空")
	exception.ThrowWebExceptionBool(content == "", 403, "配置内容不能为空")

	// 获取应用信息
	do := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionBool(do.IsNil(), 403, "应用不存在")

	client := docker.NewClient()

	// 新版本号
	newVersion := do.ConfigVer + 1

	// 创建新的 Docker Config
	configName := fmt.Sprintf("%s_config_v%d", appName, newVersion)
	labels := map[string]string{
		"owner_service": appName,
		"version":       fmt.Sprintf("%d", newVersion),
	}

	_, err := client.Config.Create(configName, []byte(content), labels)
	if err != nil {
		exception.ThrowWebExceptionf(403, "创建 Docker Config 失败: %v", err)
	}

	// 更新应用的配置版本号
	do.ConfigVer = newVersion
	err = appsRepository.UpdateApp(do)
	exception.ThrowWebExceptionError(403, err)
}
