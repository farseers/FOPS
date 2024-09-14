package ruleTimeType

// Enum 规则时间类型
type Enum int

const (
	Hour Enum = iota // 小时
	Day              // 天
)

func (e Enum) ToString() string {
	switch e {
	case Hour:
		return "小时"
	case Day:
		return "天"
	}
	return ""
}
