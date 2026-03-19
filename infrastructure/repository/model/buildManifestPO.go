package model

import (
	"time"

	"github.com/farseer-go/data"
)

// BuildManifestPO 构建清单持久化对象
type BuildManifestPO struct {
	AppName       string    `gorm:"primaryKey;size:32;not null;comment:应用名称"`
	GitName       string    `gorm:"primaryKey;size:64;not null;comment:应用或库名称"`
	BuildNumber   int       `gorm:"primaryKey;type:int;not null;default:0;comment:构建号"`
	WorkflowsName string    `gorm:"size:32;not null;default:'';comment:工作流名称"`
	DockerImage   string    `gorm:"size:128;not null;default:'';comment:镜像名称"`
	GitId         int       `gorm:"type:int;not null;default:0;comment:Git主键"`
	GitBranch     string    `gorm:"size:64;not null;default:'';comment:GIT分支"`
	GitCommitId   string    `gorm:"size:64;not null;default:'';comment:git commitId"`
	CreateAt      time.Time `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:构建时间"`
}

// CreateIndex 创建索引
func (*BuildManifestPO) CreateIndex() map[string]data.IdxField {
	return map[string]data.IdxField{
		"idx_app_name":     {IsUNIQUE: false, Fields: "app_name, create_at desc"},
		"idx_docker_image": {IsUNIQUE: false, Fields: "docker_image"},
		"idx_create_at":    {IsUNIQUE: false, Fields: "create_at desc"},
	}
}
