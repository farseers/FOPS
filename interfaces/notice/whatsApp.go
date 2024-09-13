package notice

import (
	"fmt"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/utils/http"
	"net/url"
	"strings"
)

// WhatsAppSendMsg 发送whatsApp消息  phone:apikey
func WhatsAppSendMsg(content string, apiKeys []string) {
	for _, key := range apiKeys {
		apiArray := strings.Split(key, ":")
		if len(apiArray) != 2 {
			continue
		}
		sendUrl := fmt.Sprintf("https://api.callmebot.com/whatsapp.php?phone=%s&text=%s&apikey=%s", apiArray[0], url.QueryEscape(content), apiArray[1])
		head := make(map[string]any)
		head["User-Agent"] = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36"
		res, _, _, err := http.NewClient(sendUrl).Head(head).Timeout(5000).Get()
		flog.ErrorIfExists(err)
		flog.Info("whatsapp 发送消息返回数据：" + res)
	}
}
