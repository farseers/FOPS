package eumBuildType

type Enum int

const (
	Manual Enum = iota // 手动构建
	Auto               // 自动构建
)
