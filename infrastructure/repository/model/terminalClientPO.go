package model

type TerminalClientPO struct {
	Name      string `gorm:"size:32;not null;comment:客户端名称"`
	LoginIp   string `gorm:"primaryKey;size:32;not null;comment:客户方ip"`
	LoginName string `gorm:"size:32;not null;comment:登录名"`
	LoginPwd  string `gorm:"size:32;not null;comment:登录密码"`
	LoginPort int    `gorm:"type:int;not null;comment:端口"`
}
