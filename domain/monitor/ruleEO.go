package monitor

type RuleEO struct {
	Id                   int64  // 主键
	AppId                string // 项目ID
	AppName              string // 项目名称
	Comparison           string // 比较方式 >  =  <
	KeyName              string // 监控键
	KeyValue             string // 监控键值
	Remark               string // 备注
	Enable               bool   // 是否启用
	NoticeWhatsAppApiKey string // whatsapp通知key
}

func (receiver *RuleEO) IsNull() bool {
	return receiver == nil || receiver.Id == 0 || len(receiver.AppId) == 0
}
