package model

type MonitorRulePO struct {
	Id                   int64  `gorm:"primaryKey;autoIncrement;comment:主键"`
	AppId                string `gorm:"size:32;not null;comment:项目ID"`
	AppName              string `gorm:"size:32;not null;comment:项目名称"`
	Comparison           string `gorm:"size:32;not null;comment:比较方式 >  =  <"`
	KeyName              string `gorm:"size:32;not null;comment:监控键"`
	KeyValue             string `gorm:"size:32;not null;comment:监控键值"`
	Remark               string `gorm:"size:256;not null;comment:备注"`
	NoticeWhatsAppApiKey string `gorm:"size:256;not null;comment:whatsapp通知key"`
	Enable               bool   `gorm:"size:1;not null;default:0;comment:是否启用"`
}
