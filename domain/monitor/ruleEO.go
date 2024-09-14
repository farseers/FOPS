package monitor

import (
	"fmt"
	"fops/domain/enum/ruleTimeType"
	"fops/interfaces/notice"
	"github.com/farseer-go/fs/parse"
	"strings"
	"time"
)

type RuleEO struct {
	Id                   int64             // 主键
	AppId                string            // 项目ID
	AppName              string            // 项目名称
	TimeType             ruleTimeType.Enum // 规则时间类型 0小时，1天
	StartTime            time.Time         // 开始时间
	EndTime              time.Time         // 结束时间
	Comparison           string            // 比较方式 >  =  <
	KeyName              string            // 监控键
	KeyValue             string            // 监控键值
	Remark               string            // 备注
	Enable               bool              // 是否启用
	NoticeWhatsAppApiKey string            // whatsapp通知key
}

// IsNull 判断是否为空
func (receiver *RuleEO) IsNull() bool {
	return receiver == nil || receiver.Id == 0 || len(receiver.AppId) == 0
}

// 发送whatsapp 消息
func (receiver *RuleEO) WhatsAppSendMsg(reqVal string) {
	// 时间类型判断
	var send = false
	switch receiver.TimeType {
	case ruleTimeType.Hour:
		if time.Now().Hour() >= receiver.StartTime.Hour() && time.Now().Hour() <= receiver.EndTime.Hour() {
			send = true
		}
	case ruleTimeType.Day:
		if time.Now().After(receiver.StartTime) && time.Now().Before(receiver.EndTime) {
			send = true
		}
	}
	if !receiver.IsNull() && send {
		var comparisonMsg string
		switch receiver.Comparison {
		case ">":
			if parse.ToFloat32(receiver.KeyValue) > parse.ToFloat32(reqVal) {
				comparisonMsg = "大于"
			}
		case "<":
			if parse.ToFloat32(receiver.KeyValue) < parse.ToFloat32(reqVal) {
				comparisonMsg = "小于"
			}
		case "=":
			if parse.ToFloat32(receiver.KeyValue) == parse.ToFloat32(reqVal) {
				comparisonMsg = "等于"
			}
		}
		// 发送消息
		if comparisonMsg != "" && receiver.NoticeWhatsAppApiKey != "" {
			notice.WhatsAppSendMsg(fmt.Sprintf("%s %s：%s，%s：%s", time.Now().Format("2006-01-02 15:04:05"), receiver.AppName, receiver.Remark, comparisonMsg, reqVal), strings.Split(receiver.NoticeWhatsAppApiKey, ","))
		}
	}
}
