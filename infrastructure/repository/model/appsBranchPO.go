package model

import (
	"github.com/farseer-go/fs/dateTime"
)

// AppsBranchPO 实体类
type AppsBranchPO struct {
	AppName         string            `gorm:"primaryKey;size:32;not null;comment:应用名称"`
	BranchName      string            `gorm:"primaryKey;size:64;not null;comment:分支名称"`
	CommitId        string            `gorm:"primaryKey;size:64;not null;comment:当前分支最后提交ID"`
	CommitMessage   string            `gorm:"size:256;not null;default:'';comment:提交消息"`
	CommitAt        dateTime.DateTime `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:同步时间"`
	BuildId         int64             `gorm:"not null;default:0;comment:对应的构建ID"`
	BuildSuccess    bool              `gorm:"size:1;not null;default:0;comment:是否构建成功"`
	BuildErrorCount int               `gorm:"type:int;not null;default:0;comment:构建失败次数"`
	BuildAt         dateTime.DateTime `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:构建时间"`
}
