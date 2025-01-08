package model

import (
	"fops/domain/enum/ruleTimeType"
)

type MonitorRulePO struct {
	Id          int64             `gorm:"primaryKey;autoIncrement;comment:主键"`
	AppName     string            `gorm:"size:512;not null;comment:应用名称"`
	TimeType    ruleTimeType.Enum `gorm:"not null;comment:规则时间类型 0小时，1天"`
	StartDate   string            `gorm:"size:32;not null;comment:开始小时"`
	EndDate     string            `gorm:"size:32;not null;comment:结束小时"`
	StartDay    string            `gorm:"size:32;not null;comment:开始天"`
	EndDay      string            `gorm:"size:32;not null;comment:结束天"`
	Comparison  string            `gorm:"size:32;not null;comment:比较方式 >  =  <"`
	KeyName     string            `gorm:"size:32;not null;comment:监控键"`
	KeyValue    string            `gorm:"size:32;not null;comment:监控键值"`
	Remark      string            `gorm:"size:256;comment:备注"`
	TipTemplate string            `gorm:"size:256;comment:提示模版"`
	NoticeIds   []int             `gorm:"json;comment:NoticeIds"`
	Enable      bool              `gorm:"size:1;not null;default:0;comment:是否启用"`
}
