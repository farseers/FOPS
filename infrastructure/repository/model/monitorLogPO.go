package model

import (
	"github.com/farseer-go/collections"
	"time"
)

type MonitorLogPO struct {
	Id       int64                                  `gorm:"primaryKey;autoIncrement;comment:主键"`
	AppId    string                                 `gorm:"size:32;not null;comment:项目ID"`
	AppName  string                                 `gorm:"size:256;not null;comment:项目名称"`
	Keys     collections.Dictionary[string, string] `gorm:"size:2048;json;not null;comment:监控键值对"`
	CreateAt time.Time                              `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:开始时间"`
}
