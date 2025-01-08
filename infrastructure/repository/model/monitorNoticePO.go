package model

import "fops/domain/enum/noticeType"

type MonitorNoticePO struct {
	Id         int64           `gorm:"primaryKey;autoIncrement;comment:主键"`
	NoticeType noticeType.Enum `gorm:"type:tinyint;not null;default:0;comment:通知类型：0=whatsapp"`
	Name       string          `gorm:"size:32;not null;comment:名称"`
	Email      string          `gorm:"size:128;not null;comment:Email"`
	Phone      string          `gorm:"size:32;not null;comment:号码"`
	ApiKey     string          `gorm:"size:256;not null;comment:接口Key"`
	Remark     string          `gorm:"size:256;not null;comment:备注"`
	Enable     bool            `gorm:"size:1;not null;default:0;comment:是否启用"`
}
