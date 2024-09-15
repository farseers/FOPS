package monitor

type NoticeEO struct {
	Id     int64  // 主键
	Email  string // 邮箱
	Phone  string // 号码
	ApiKey string // 接口Key
	Remark string // 备注
	Enable bool   // 是否启用
}
