package monitor

import (
	"fmt"
	"fops/domain/enum/noticeType"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/utils/http"
	"net/url"
)

type NoticeEO struct {
	Id         int64           // 主键
	NoticeType noticeType.Enum // 通知类型
	Name       string          // 名称
	Email      string          // 邮箱
	Phone      string          // 号码
	ApiKey     string          // 接口Key
	Remark     string          // 备注
	Enable     bool            // 是否启用
}

// Notice 通知
func (receiver *NoticeEO) Notice(content string) {
	var (
		body       string
		statusCode int
		err        error
		head       = map[string]any{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"}
	)

	switch receiver.NoticeType {
	case noticeType.WhatsApp: // whatsApp
		sendUrl := fmt.Sprintf("https://api.callmebot.com/whatsapp.php?phone=%s&apikey=%s&text=%s", receiver.Phone, receiver.ApiKey, url.QueryEscape(content))
		body, statusCode, _, err = http.RequestProxyConfigure("GET", sendUrl, head, nil, "", 5000)
	case noticeType.Telegram: // Telegram
		sendUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendmessage?chat_id=%s&text=%s", receiver.ApiKey, receiver.Phone, url.QueryEscape(content))
		body, statusCode, _, err = http.RequestProxyConfigure("GET", sendUrl, head, nil, "", 5000)
	case noticeType.Log:
		statusCode = 200
		flog.Infof("消息通知：%s", content)
	}

	if err != nil {
		flog.Warningf("发送告警通知异常：%s", err.Error())
		return
	}
	if statusCode != 200 {
		flog.Warningf("【%s】发送告警通知失败：statusCode = %d %s", receiver.NoticeType.ToString(), statusCode, body)
	}
}
