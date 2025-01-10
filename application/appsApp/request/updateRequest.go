package request

import "github.com/farseer-go/collections"

type UpdateRequest struct {
	ClusterDockerImage string                  // 集群镜像名称
	AppName            string                  `validate:"required" label:"应用名称"` // 应用名称
	AppGit             int64                   // 应用的源代码
	FrameworkGits      collections.List[int64] // 依赖的框架源代码
	DockerfilePath     string                  // DockerfilePath路径
	HealthInstance     int                     // 健康的实例数量
	DockerReplicas     int                     // 副本数量
	DockerNodeRole     string                  // 容器节点角色 manager or worker
	AdditionalScripts  string                  // 首次创建应用时附加脚本
	WorkflowsYmlPath   string                  // 工作流定义的路径（默认：/.fops/workflows/build.yml）
	LimitCpus          float64                 // Cpu核数限制
	LimitMemory        string                  // 内存限制
	UTWorkflowsName    string                  // UT工作流名称（文件的名称）
}
