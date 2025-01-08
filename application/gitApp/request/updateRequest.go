package request

type UpdateRequest struct {
	Id       int64  // 主键
	Name     string // Git名称
	Hub      string // git地址
	Branch   string // Git分支
	UserName string // 账户名称
	UserPwd  string // 账户密码
	Path     string // 存储目录
	IsApp    bool   // 是否为应用
}
