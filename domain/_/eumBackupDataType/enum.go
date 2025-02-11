package eumBackupDataType

// Enum 备份数据类型
type Enum int

const (
	Mysql      Enum = iota // Mysql
	Clickhouse             // Clickhouse
)
