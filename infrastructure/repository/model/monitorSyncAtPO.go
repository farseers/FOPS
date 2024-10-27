package model

import (
	"time"
)

type MonitorSyncAtPO struct {
	AppName string    `gorm:"size:64;not null;comment:应用名称"`
	SyncAt  time.Time `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:同步时间"`
}
