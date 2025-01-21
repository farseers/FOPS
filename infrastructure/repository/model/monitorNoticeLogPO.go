package model

import (
	"fops/domain/enum/noticeType"
	"time"
)

type MonitorNoticeLogPO struct {
	Id         int64           `gorm:"primaryKey;autoIncrement;comment:主键"`
	AppName    string          `gorm:"size:64;not null;comment:项目名称"`
	NoticeId   int64           `gorm:"not null;comment:通知Id"`
	NoticeName string          `gorm:"size:32;not null;comment:通知人"`
	NoticeType noticeType.Enum `gorm:"type:tinyint;not null;default:0;comment:通知类型：0=whatsapp"`
	NoticeMsg  string          `gorm:"size:4096;not null;comment:通知消息"`
	NoticeAt   time.Time       `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:通知时间"`
	IsRead     bool            `gorm:"size:1;not null;default:0;comment:是否已读"`
}
