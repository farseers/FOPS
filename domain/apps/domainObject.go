package apps

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
)

// DomainObject 聚合
type DomainObject struct {
	AppName           string                  // 应用名称（链路追踪）
	DockerVer         int                     // 镜像版本
	ClusterVer        map[int64]*ClusterVerVO // 集群版本
	AppGit            int64                   // 应用的源代码
	FrameworkGits     collections.List[int64] // 依赖的框架源代码
	DockerImage       string                  // 仓库镜像名称
	DockerfilePath    string                  // Dockerfile路径
	DockerInstances   int                     // 运行的实例数量
	DockerInspect     []DockerInspectVO       // 运行的实例详情
	DockerReplicas    int                     // 副本数量
	DockerNodeRole    string                  // 容器节点角色 manager or worker or global
	AdditionalScripts string                  // 首次创建应用时附加脚本
	LimitCpus         float64                 // Cpu核数限制
	LimitMemory       string                  // 内存限制
	IsSys             bool                    // 是否系统应用
}

func (receiver *DomainObject) IsNil() bool {
	return receiver.AppName == ""
}

// ClusterVerVO 集群镜像版本及部署时间
type ClusterVerVO struct {
	ClusterId       int64             // 集群ID
	DockerVer       int               // 集群镜像版本
	DockerImage     string            // 集群镜像名称
	DeploySuccessAt dateTime.DateTime // 上次部署成功时间
	BuildSuccessId  int64             // 上次部署成功的构建ID
	DeployFailAt    dateTime.DateTime // 上次部署失败时间
	BuildFailId     int64             // 上次部署失败的构建ID
}

// UpdateBuildVer 当构建失败时，记录失败时间、失败时的构建ID
func (receiver *DomainObject) UpdateBuildVer(isSuccess bool, clusterId int64, buildId int64) {
	if receiver.ClusterVer == nil {
		receiver.ClusterVer = make(map[int64]*ClusterVerVO)
	}
	if _, ok := receiver.ClusterVer[clusterId]; !ok {
		receiver.ClusterVer[clusterId] = &ClusterVerVO{
			ClusterId: clusterId,
		}
	}

	// 当构建成功时，记录发布时间、发布时的构建ID
	if isSuccess {
		receiver.ClusterVer[clusterId].DockerVer = receiver.DockerVer
		receiver.ClusterVer[clusterId].DockerImage = receiver.DockerImage
		receiver.ClusterVer[clusterId].DeploySuccessAt = dateTime.Now()
		receiver.ClusterVer[clusterId].BuildSuccessId = buildId
	} else // 当构建失败时，记录失败时间、失败时的构建ID
	{
		receiver.ClusterVer[clusterId].DeployFailAt = dateTime.Now()
		receiver.ClusterVer[clusterId].BuildFailId = buildId
	}
}

// GetWorkflowsRoot 获取工作流目录 如："/var/lib/fops/workflows/fops/"
func (receiver *DomainObject) GetWorkflowsRoot() string {
	return WorkflowsRoot + receiver.AppName + "/"
}

// GetWorkflowsDir 获取工作流目录 如："/var/lib/fops/workflows/fops/.fops/workflows/"
func (receiver *DomainObject) GetWorkflowsDir() string {
	return WorkflowsRoot + receiver.AppName + "/.fops/workflows/"
}

// GetCurClusterDockerImage 获取当前集群的镜像名称
func (receiver *DomainObject) GetCurClusterDockerImage(clusterId int64) string {
	if cur, ok := receiver.ClusterVer[clusterId]; ok {
		return cur.DockerImage
	}
	return ""
}

func (receiver *DomainObject) InitCluster(clusterId int64) {
	if receiver.ClusterVer == nil {
		receiver.ClusterVer = make(map[int64]*ClusterVerVO)
	}
	if receiver.ClusterVer[clusterId] == nil {
		receiver.ClusterVer[clusterId] = &ClusterVerVO{
			ClusterId: clusterId,
		}
	}
}
