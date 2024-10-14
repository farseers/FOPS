package monitor

import (
	"fops/domain/enum/noticeType"
	"time"
)

type NoticeLogEO struct {
	Id         int64           // 主键
	AppName    string          // 项目名称
	NoticeId   int64           // 通知Id
	NoticeName string          // 通知人
	NoticeType noticeType.Enum // 0 whatsapp
	NoticeMsg  string          // 通知消息
	NoticeAt   time.Time       // 通知时间
	IsRead     bool            // 是否已读
}

// NewLog 添加新日志
func NewLog(appName string, noticeId int64, noticeName string, notType noticeType.Enum, noticeMsg string) NoticeLogEO {
	return NoticeLogEO{
		AppName:    appName,
		NoticeId:   noticeId,
		NoticeName: noticeName,
		NoticeType: notType,
		NoticeMsg:  noticeMsg,
		NoticeAt:   time.Now(),
		IsRead:     false,
	}
}
