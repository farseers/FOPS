package model

import (
	"fops/domain/apps"
	"github.com/farseer-go/collections"
)

type ClusterNodePO struct {
	NodeName      string                               `gorm:"primaryKey;comment:节点名称"`
	Status        string                               `gorm:"size:64;not null;comment:主机状态"`
	Availability  string                               `gorm:"size:64;not null;comment:节点状态"`
	IsMaster      bool                                 `gorm:"size:1;not null;default:0;comment:是否为主节点"`
	EngineVersion string                               `gorm:"size:64;not null;comment:引擎版本"`
	IP            string                               `gorm:"size:64;not null;comment:节点IP"`
	OS            string                               `gorm:"size:64;not null;comment:操作系统"`
	Architecture  string                               `gorm:"size:64;not null;comment:架构"`
	CPUs          string                               `gorm:"size:64;not null;comment:CPU核心数"`
	Memory        string                               `gorm:"size:64;not null;comment:内存"`
	Label         collections.List[apps.DockerLabelVO] `gorm:"size:2048;json;not null;comment:标签"`
}
