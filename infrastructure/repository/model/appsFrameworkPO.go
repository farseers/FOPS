package model

// AppsFrameworkPO 应用与框架关系表
type AppsFrameworkPO struct {
	AppName     string `gorm:"primaryKey;size:32;not null;comment:应用名称"`
	FrameworkId int64  `gorm:"primaryKey;not null;comment:框架ID"`
	CommitId    string `gorm:"size:64;not null;default:'';comment:框架提交ID"`
}
