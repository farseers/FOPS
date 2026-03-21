package monitor

import (
	"fmt"
	"fops/domain/enum/noticeType"
	"net/url"
	"sync"
	"time"

	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/utils/http"
)

// telegramRateLimit 以 ApiKey+Phone 为 key，记录限流截止时间
var telegramRateLimit sync.Map

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
		// 短时间内发送了大量消息，触发限流，忽略后续1分钟内所有消息发送
		rateLimitKey := receiver.ApiKey + receiver.Phone
		if val, ok := telegramRateLimit.Load(rateLimitKey); ok && time.Now().Before(val.(time.Time)) {
			//flog.Warningf("【%s】限流冷却中，忽略消息发送，解除时间：%s", receiver.NoticeType.ToString(), val.(time.Time).Format(time.DateTime))
			return
		}
		sendUrl := fmt.Sprintf("https://api.telegram.org/bot%s/sendmessage?chat_id=%s&text=%s", receiver.ApiKey, receiver.Phone, url.QueryEscape(content))
		body, statusCode, _, err = http.RequestProxyConfigure("GET", sendUrl, nil, nil, "", 5000)
		if statusCode == 429 {
			rateLimitEnd := time.Now().Add(time.Minute)
			telegramRateLimit.Store(rateLimitKey, rateLimitEnd)
			flog.Warningf("【%s】触发限流（429），将在1分钟内忽略所有消息发送，解除时间：%s", receiver.NoticeType.ToString(), rateLimitEnd.Format(time.DateTime))
			return
		}
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
