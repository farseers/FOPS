package eumBackupStoreType

type Enum int

const (
	OSS            Enum = iota // aliyun oss
	LocalDirectory             // 本地目录
)
