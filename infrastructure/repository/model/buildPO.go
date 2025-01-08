package model

import (
	"fops/domain/_/eumBuildStatus"
	"fops/domain/_/eumBuildType"
	"fops/domain/apps"
	"time"

	"github.com/farseer-go/data"
)

// BuildPO 实体类
type BuildPO struct {
	Id            int64               `gorm:"primaryKey;comment:主键"`
	AppName       string              `gorm:"size:32;not null;comment:应用名称"`
	ClusterId     int64               `gorm:"not null;default:0;comment:集群信息"`
	BuildNumber   int                 `gorm:"type:int;not null;default:0;comment:构建号"`
	Status        eumBuildStatus.Enum `gorm:"type:tinyint;not null;default:0;comment:状态：0=未开始，1=构建中，2=完成"`
	BuildType     eumBuildType.Enum   `gorm:"type:tinyint;not null;default:0;comment:构建类型：0=应用，1=UT"`
	IsSuccess     bool                `gorm:"size:1;not null;default:0;comment:是否成功"`
	CreateAt      time.Time           `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:开始时间"`
	FinishAt      time.Time           `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:完成时间"`
	BuildServerId int64               `gorm:"not null;default:0;comment:构建的服务端id"`
	WorkflowsName string              `gorm:"size:32;not null;default:'';comment:工作流名称（文件的名称）"`
	BranchName    string              `gorm:"size:32;not null;default:'';comment:分支名称"`
	DockerImage   string              `gorm:"size:64;not null;default:'';comment:Docker镜像"`
	Env           apps.EnvVO          `gorm:"type:text;json;not null;comment:环境变量"`
}

// 创建索引
func (*BuildPO) CreateIndex() map[string]data.IdxField {
	return map[string]data.IdxField{
		"idx_app_name": {false, "app_name, build_number, id"},
		"idx_status":   {false, "status, build_type"},
	}
}
