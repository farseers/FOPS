package model

import (
	"fops/domain/apps"
	"github.com/farseer-go/collections"
)

// AppsPO 实体类
type AppsPO struct {
	AppName           string                       `gorm:"primaryKey;size:32;not null;comment:应用名称（链路追踪）"`
	DockerVer         int                          `gorm:"type:int;not null;comment:镜像版本"`
	ShellScript       string                       `gorm:"type:mediumtext;not null;comment:Shell脚本"`
	ClusterVer        map[int64]*apps.ClusterVerVO `gorm:"size:256;json;not null;comment:集群版本"`
	AppGit            int64                        `gorm:"not null;default:0;comment:应用的源代码"`
	FrameworkGits     collections.List[int64]      `gorm:"size:64;json;not null;comment:依赖的框架源代码"`
	DockerImage       string                       `gorm:"type:text;not null;comment:仓库镜像名称"`
	Dockerfile        string                       `gorm:"type:text;not null;comment:Dockerfile内容"`
	DockerfilePath    string                       `gorm:"size:256;not null;comment:Dockerfile路径"`
	ActiveInstance    []apps.ActiveInstanceEO      `gorm:"size:1024;json;not null;comment:正在运行的实例"`
	DockerReplicas    int                          `gorm:"type:int;not null;comment:副本数量"`
	DockerNodeRole    string                       `gorm:"size:256;not null;comment:容器节点角色"`
	AdditionalScripts string                       `gorm:"type:text;not null;comment:首次创建应用时附加脚本"`
}
