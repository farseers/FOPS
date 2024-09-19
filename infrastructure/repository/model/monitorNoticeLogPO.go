package model

import (
	"fops/domain/enum/noticeType"
	"time"
)

type MonitorNoticeLogPO struct {
	Id         int64           `gorm:"primaryKey;autoIncrement;comment:主键"`
	AppId      string          `gorm:"size:32;not null;comment:项目ID"`
	AppName    string          `gorm:"size:32;not null;comment:项目名称"`
	NoticeId   int64           `gorm:"not null;comment:通知Id"`
	NoticeType noticeType.Enum `gorm:"type:tinyint;not null;default:0;comment:通知类型：0=whatsapp"`
	NoticeMsg  string          `gorm:"size:256;not null;comment:通知消息"`
	NoticeAt   time.Time       `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:通知时间"`
}
