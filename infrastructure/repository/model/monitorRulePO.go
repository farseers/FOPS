package model

import (
	"fops/domain/enum/ruleTimeType"
	"time"
)

type MonitorRulePO struct {
	Id         int64             `gorm:"primaryKey;autoIncrement;comment:主键"`
	AppId      string            `gorm:"size:32;not null;comment:项目ID"`
	AppName    string            `gorm:"size:32;not null;comment:项目名称"`
	TimeType   ruleTimeType.Enum `gorm:"not null;comment:规则时间类型 0小时，1天"`
	StartTime  time.Time         `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:开始时间"`
	EndTime    time.Time         `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:结束时间"`
	Comparison string            `gorm:"size:32;not null;comment:比较方式 >  =  <"`
	KeyName    string            `gorm:"size:32;not null;comment:监控键"`
	KeyValue   string            `gorm:"size:32;not null;comment:监控键值"`
	Remark     string            `gorm:"size:256;not null;comment:备注"`
	NoticeIds  []int             `gorm:"json;not null;comment:NoticeIds"`
	Enable     bool              `gorm:"size:1;not null;default:0;comment:是否启用"`
}
