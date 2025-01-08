package model

// ConfigurePO 配置中心
type ConfigurePO struct {
	AppName string `gorm:"primaryKey;size:32;not null;comment:应用名称"`
	Key     string `gorm:"primaryKey;not null;default:'';comment:配置KEY"`
	Ver     int    `gorm:"primaryKey;type:int;not null;default:0;comment:版本"`
	Value   string `gorm:"size:1024;not null;comment:配置VALUE"`
}
