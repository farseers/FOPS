package eumRuntimeEnv

type Enum int

const (
	Dev        Enum = iota // Dev 开发环境
	Test                   // Test 测试环境
	PreRelease             // PreRelease 预发布
	Prod                   // Prod 生产环境
)
