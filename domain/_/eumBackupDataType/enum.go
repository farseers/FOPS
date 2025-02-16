package eumBackupDataType

// Enum 备份数据类型
type Enum int

const (
	Mysql      Enum = iota // Mysql
	Clickhouse             // Clickhouse
)

func (receiver Enum) ToString() string {
	switch receiver {
	case Mysql:
		return "mysql"
	case Clickhouse:
		return "clickhouse"
	}
	return ""
}
