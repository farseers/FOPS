package model

import (
	"fops/domain/_/eumBackupDataType"
	"fops/domain/_/eumBackupStoreType"
	"time"
)

// BackupDataPO 实体类
type BackupDataPO struct {
	Id             int64                   `gorm:"primaryKey;comment:订单ID"`
	BackupDataType eumBackupDataType.Enum  `gorm:"type:tinyint;not null;default:0;comment:备份数据类型"`
	Host           string                  `gorm:"size:32;not null;comment:主机"`
	Port           int                     `gorm:"type:int;not null;comment:端口"`
	Username       string                  `gorm:"size:32;not null;comment:用户名"`
	Password       string                  `gorm:"size:32;not null;comment:密码"`
	Database       []string                `gorm:"size:64;json;not null;comment:数据库"`
	LastBackupAt   time.Time               `gorm:"type:timestamp;size:6;not null;comment:上次备份时间"`
	Cron           string                  `gorm:"size:32;not null;comment:备份间隔"`
	StoreType      eumBackupStoreType.Enum `gorm:"type:tinyint;not null;default:0;comment:备份存储类型"`
	StoreConfig    string                  `gorm:"size:32;not null;comment:备份存储配置"`
}
