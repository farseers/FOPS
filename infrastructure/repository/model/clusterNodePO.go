package model

import (
	"fops/domain/apps"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
)

type ClusterNodePO struct {
	NodeName           string                               `gorm:"primaryKey;comment:节点名称"`
	Status             string                               `gorm:"size:64;not null;comment:主机状态"`
	Availability       string                               `gorm:"size:64;not null;comment:节点状态"`
	IsMaster           bool                                 `gorm:"size:1;not null;default:0;comment:是否为主节点"`
	IsHealth           bool                                 `gorm:"size:1;not null;default:0;comment:应用是否健康"`
	EngineVersion      string                               `gorm:"size:64;not null;comment:引擎版本"`
	IP                 string                               `gorm:"size:64;not null;comment:节点IP"`
	AgentIP            string                               `gorm:"size:64;not null;comment:代理容器IP"`
	OS                 string                               `gorm:"size:64;not null;comment:操作系统"`
	Architecture       string                               `gorm:"size:64;not null;comment:架构"`
	CPUs               string                               `gorm:"size:64;not null;comment:CPU核心数"`
	Memory             string                               `gorm:"size:64;not null;comment:内存"`
	CpuUsagePercent    float64                              `gorm:"type:decimal(6,1);size:64;not null;comment:CPU使用百分比"`
	MemoryUsagePercent float64                              `gorm:"type:decimal(6,1);size:64;not null;comment:内存使用百分比"`
	MemoryUsage        float64                              `gorm:"size:64;not null;comment:内存已使用（MB）"`
	Disk               uint64                               `gorm:"size:64;not null;comment:硬盘总容量（GB）"`
	DiskUsagePercent   float64                              `gorm:"type:decimal(6,1);size:64;not null;comment:硬盘使用百分比"`
	DiskUsage          float64                              `gorm:"size:64;not null;comment:硬盘已用空间（GB）"`
	Label              collections.List[apps.DockerLabelVO] `gorm:"size:2048;json;not null;comment:标签"`
	UpdateAt           dateTime.DateTime                    `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:更新时间"`
}
