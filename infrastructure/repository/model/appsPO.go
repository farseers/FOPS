package model

import (
	"fops/domain/apps"

	"github.com/farseer-go/collections"
)

// AppsPO 实体类
type AppsPO struct {
	AppName           string                                 `gorm:"primaryKey;size:32;not null;comment:应用名称（链路追踪）"`
	ClusterVer        map[int64]apps.ClusterVerVO            `gorm:"size:2048;json;not null;comment:集群版本"`
	AppGit            int64                                  `gorm:"not null;default:0;comment:应用的源代码"`
	FrameworkGits     collections.List[int64]                `gorm:"size:64;json;not null;comment:依赖的框架源代码"`
	DockerfilePath    string                                 `gorm:"size:256;not null;default:'';comment:Dockerfile路径"`
	DockerInstances   int                                    `gorm:"type:int;not null;default:0;comment:运行的实例数量"`
	DockerInspect     collections.List[apps.DockerInspectVO] `gorm:"size:2048;json;not null;comment:运行的实例详情"`
	DockerReplicas    int                                    `gorm:"type:int;not null;default:0;comment:副本数量"`
	DockerNodeRole    string                                 `gorm:"size:256;not null;default:'';comment:容器节点角色"`
	AdditionalScripts string                                 `gorm:"type:text;not null;comment:首次创建应用时附加脚本"`
	UTWorkflowsName   string                                 `gorm:"size:32;not null;default:'';comment:UT工作流名称（文件的名称）"`
	LimitCpus         float64                                `gorm:"type:decimal(6,3);not null;default:0;comment:Cpu核数限制"`
	LimitMemory       string                                 `gorm:"not null;default:'';comment:内存限制"`
	IsSys             bool                                   `gorm:"size:1;not null;default:0;comment:是否系统应用"`
}
