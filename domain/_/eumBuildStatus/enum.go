package eumBuildStatus

type Enum int

const (
	None     Enum = iota // None 未开始
	Building             // Building 构建中
	Finish               // Finish 完成
	Cancel               // 取消
)
