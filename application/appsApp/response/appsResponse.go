package response

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
)

type ClusterVerVO struct {
	ClusterId       int64             // 集群ID
	ClusterName     string            // 集群名称
	DockerImage     string            // 集群镜像名称
	DeploySuccessAt dateTime.DateTime // 上次部署成功时间
}

type AppsResponse struct {
	AppName            string                         // 应用名称（链路追踪）
	DockerVer          int                            // 镜像版本
	ClusterVer         collections.List[ClusterVerVO] // 集群版本
	LocalClusterVer    ClusterVerVO                   // 集群版本
	DockerImage        string                         // 仓库镜像名称
	AppGit             int64                          // 应用的git仓库
	FrameworkGits      collections.List[int64]        // 依赖的框架源代码
	IsHealth           bool                           // 应用是否健康
	DockerInstances    int                            // 运行的实例数量
	DockerReplicas     int                            // 副本数量
	DockerNodeRole     string                         // 容器节点角色 manager or worker or global
	LogErrorCount      int                            // 日志错误数量
	LogWaringCount     int                            // 日志警告数量
	TaskFailCount      int                            // 任务失败数量
	TaskSuccessCount   int                            // 任务成功数量
	WorkflowsNames     []string                       // 工作流名称
	AdditionalScripts  string                         // 首次创建应用时附加脚本
	DockerfilePath     string                         // DockerfilePath路径
	LimitCpus          float64                        // Cpu核数限制
	LimitMemory        string                         // 内存限制
	CpuUsagePercent    float64                        // CPU使用百分比
	MemoryUsagePercent float64                        // 内存使用百分比
	MemoryUsage        uint64                         // 内存已使用（MB）
	MemoryLimit        uint64                         // 内存限制（MB）
}
