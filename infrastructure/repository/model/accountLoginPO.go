package model

import "time"

type AccountLoginPO struct {
	Id        uint64    `gorm:"primaryKey;autoIncrement;comment:主键"`
	LoginName string    `gorm:"size:32;not null;default:'';comment:登陆名称"`
	LoginPwd  string    `gorm:"size:32;not null;default:'';comment:登录密码"`
	LoginSalt string    `gorm:"size:32;not null;default:'';comment:盐"`
	CreateAt  time.Time `gorm:"type:timestamp;size:6;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
}
