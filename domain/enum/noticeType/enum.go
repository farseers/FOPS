package noticeType

import "github.com/farseer-go/collections"

// Enum 通知类型
type Enum int

const (
	WhatsApp Enum = iota // whatsApp
	Telegram             // Telegram
	Log      = -1        // 仅打印日志
)

func (e Enum) ToString() string {
	switch e {
	case WhatsApp:
		return "whatsApp"
	case Telegram:
		return "Telegram"
	case Log:
		return "Log"
	}
	return ""
}

func ToList() collections.List[Enum] {
	return collections.NewList[Enum](WhatsApp, Telegram, Log)
}
