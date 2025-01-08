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
	PullWorkflows(ctx context.Context, gitPath, branch string, gitRemote string, progress chan string) bool
	// 获取远程分支和提交的CommitId、Message
	GetRemoteBranch(ctx context.Context, gitPath string) collections.List[RemoteBranchVO]
}

// DockerLabelVO 标签
type DockerLabelVO struct {
	Name  string // 标签名称
	Value string // 标签值
}
