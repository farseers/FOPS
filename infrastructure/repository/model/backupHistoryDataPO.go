package model

import (
	"fops/domain/_/eumBackupStoreType"
	"time"

	"github.com/farseer-go/data"
)

// BackupHistoryDataPO 备份文件列表
type BackupHistoryDataPO struct {
	BackupId  string                  `gorm:"primaryKey;size:128;not null;comment:备份计划的ID"`
	FileName  string                  `gorm:"primaryKey;size:128;not null;comment:文件名"`
	StoreType eumBackupStoreType.Enum `gorm:"type:tinyint;not null;default:0;comment:备份存储类型"`
	CreateAt  time.Time               `gorm:"type:timestamp;size:6;not null;comment:备份时间"`
	Size      int64                   `gorm:"not null;comment:备份文件大小（KB）"`
}

// 创建索引
func (*BackupHistoryDataPO) CreateIndex() map[string]data.IdxField {
	return map[string]data.IdxField{}
}
