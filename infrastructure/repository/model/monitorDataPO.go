package model

import (
	"github.com/farseer-go/fs/dateTime"
)

type MonitorDataPO struct {
	AppName  string            `gorm:"not null;default:'';comment:项目名称"`
	Key      string            `gorm:"not null;default:'';comment:监控key"`
	Value    string            `gorm:"not null;default:'';comment:监控value"`
	CreateAt dateTime.DateTime `gorm:"type:DateTime64(3);not null;comment:创建时间"`
}
