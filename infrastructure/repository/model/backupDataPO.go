package model

import (
	"fops/domain/_/eumBackupDataType"
	"fops/domain/_/eumBackupStoreType"
	"time"

	"github.com/farseer-go/data"
)

// BackupDataPO 实体类
type BackupDataPO struct {
	Id             string                  `gorm:"primaryKey;size:128;comment:ID"`
	BackupDataType eumBackupDataType.Enum  `gorm:"type:tinyint;not null;default:0;comment:备份数据类型"`
	Host           string                  `gorm:"size:32;not null;comment:主机"`
	Port           int                     `gorm:"type:int;not null;comment:端口"`
	Username       string                  `gorm:"size:32;not null;comment:用户名"`
	Password       string                  `gorm:"size:128;not null;comment:密码"`
	Database       []string                `gorm:"size:1024;json;not null;comment:数据库"`
	LastBackupAt   time.Time               `gorm:"type:timestamp;size:6;not null;comment:上次备份时间"`
	NextBackupAt   time.Time               `gorm:"type:timestamp;size:6;not null;comment:下次备份时间"`
	Cron           string                  `gorm:"size:32;not null;comment:备份间隔"`
	StoreType      eumBackupStoreType.Enum `gorm:"type:tinyint;not null;default:0;comment:备份存储类型"`
	StoreConfig    string                  `gorm:"size:2048;not null;comment:备份存储配置"`
}

// 创建索引
func (*BackupDataPO) CreateIndex() map[string]data.IdxField {
	return map[string]data.IdxField{
		"next_backup_at": {false, "next_backup_at"},
	}
}
