package request

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/utils/system"
)

type Request struct {
	Host                system.Resource                        // 主机资源
	IsDockerMaster      bool                                   // 是否是Docker主节点
	DockerEngineVersion string                                 // Docker引擎版本
	Dockers             collections.List[docker.DockerStatsVO] // Docker容器资源
	Availability        string                                 // Docker节点状态
	Label               collections.List[docker.DockerLabelVO] // Docker节点标签
	Role                string                                 // 节点角色   manager worker
}
