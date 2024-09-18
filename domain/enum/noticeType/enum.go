package noticeType

// Enum 通知类型
type Enum int

const (
	WhatsApp Enum = iota // whatsApp
)

func (e Enum) ToString() string {
	switch e {
	case WhatsApp:
		return "whatsApp"
	}
	return ""
}
