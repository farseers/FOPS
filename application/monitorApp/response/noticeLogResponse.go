package response

import (
	"fops/domain/enum/noticeType"
)

type NoticeLogResponse struct {
	Id         int64           // 主键
	AppName    string          // 项目名称
	NoticeId   int64           // 通知Id
	NoticeName string          // 通知人
	NoticeType noticeType.Enum // 0 whatsapp
	NoticeMsg  string          // 通知消息
	NoticeAt   string          // 通知时间
}
