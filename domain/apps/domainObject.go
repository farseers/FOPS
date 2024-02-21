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
	ActiveInstance    []ActiveInstanceEO      // 正在运行的实例
	DockerReplicas    int                     // 副本数量
	DockerNodeRole    string                  // 容器节点角色 manager or worker
	AdditionalScripts string                  // 首次创建应用时附加脚本
	WorkflowsYmlPath  string                  // 工作流定义的路径（默认：/.fops/workflows/build.yml）
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
