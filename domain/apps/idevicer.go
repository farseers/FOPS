package apps

import (
	"context"
	"fops/domain/_/eumK8SControllers"
	"fops/domain/cluster"
	"github.com/farseer-go/collections"
)

type IDockerDevice interface {
	// GetDockerHub 取得dockerHub
	GetDockerHub(dockerHubAddress string) string
	// GetDockerImage 生成镜像名称，如hub.fsgit.com/fops:1
	GetDockerImage(dockerHubAddress string, projectName string, buildNumber int) string
	// Run 运行容器
	Run(name string, network string, dockerImage string, args []string, useRm bool, env EnvVO, progress chan string, ctx context.Context) bool
	// Execute 执行容器命令
	Execute(name string, execCmd string, env map[string]string, progress chan string, ctx context.Context) bool
	// Copy 复制文件到容器内
	Copy(dockerName string, sourceFile, destFile string, env EnvVO, progress chan string, ctx context.Context) bool
	// ExistsDocker 判断是否有容器在运行
	ExistsDocker(dockerName string) bool
	// Kill 停止容器
	Kill(dockerName string)
	// Remove 移除容器
	Remove(dockerName string)
	// ClearImages 清除镜像
	ClearImages(progress chan string) bool
	// GetVersion 获取版本
	GetVersion() string
	// Login 登陆镜像仓库
	Login(dockerHub string, loginName string, loginPwd string, progress chan string) bool
	// Pull 拉取镜像
	Pull(image string, progress chan string)
	// Logs 获取日志
	Logs(appName string, tailCount int) collections.List[string]
}

type IDockerSwarmDevice interface {
	// DeleteService 删除容器服务
	DeleteService(appName string, progress chan string) bool
	// SetImages 更新镜像版本
	SetImages(cluster cluster.DomainObject, appName string, dockerImages string, progress chan string) bool
	// SetImagesAndReplicas 更新镜像版本和副本数量
	SetImagesAndReplicas(cluster cluster.DomainObject, appName string, dockerImages string, dockerReplicas int, progress chan string) bool
	// SetReplicas 更新副本数量
	SetReplicas(cluster cluster.DomainObject, appName string, dockerReplicas int, progress chan string) bool
	// Restart 重启容器
	Restart(cluster cluster.DomainObject, appName string, progress chan string) bool
	ExistsDocker(appName string) bool
	// CreateService 创建服务
	CreateService(appName, dockerNodeRole, additionalScripts, dockerNetwork string, dockerReplicas int, dockerImages string, progress chan string, ctx context.Context) bool
	// Logs 获取日志
	Logs(appName string, tailCount int) collections.List[string]
	// ServiceList 获取所有Service
	ServiceList() collections.List[string]
}

type IKubectlDevice interface {
	// GetConfigFile 获取存储k8s Config的路径
	GetConfigFile(clusterName string) string
	// CreateConfigFile 创建用于K8S远程管理的配置文件
	CreateConfigFile(clusterName string, clusterConfig string) string
	// SetYaml 生成yaml文件，并执行kubectl apply命令
	SetYaml(clusterName string, projectName string, yamlContent string, progress chan string, ctx context.Context) bool
	// SetImages 更新k8s的镜像版本
	SetImages(cluster cluster.DomainObject, projectName string, dockerImages string, k8SControllersType eumK8SControllers.Enum, progress chan string, ctx context.Context) bool
	// SetImagesByClusterName 更新k8s的镜像版本
	SetImagesByClusterName(clusterName string, clusterConfig string, projectName string, dockerImages string, k8SControllersType eumK8SControllers.Enum, progress chan string, ctx context.Context) bool
}

type IGitDevice interface {
	// PullWorkflows 拉取工作流
	PullWorkflows(gitPath, branch string, gitRemote string, progress chan string) bool
}
