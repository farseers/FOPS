package model

import (
	"github.com/farseer-go/fs/dateTime"
)

type MonitorDataPO struct {
	AppId    string            `gorm:"size:32;not null;comment:项目ID"`
	AppName  string            `gorm:"size:256;not null;comment:项目名称"`
	Key      string            `gorm:"size:256;not null;comment:监控key"`
	Value    string            `gorm:"size:256;not null;comment:监控value"`
	CreateAt dateTime.DateTime `gorm:"type:DateTime64(3);not null;comment:创建时间"`
}
