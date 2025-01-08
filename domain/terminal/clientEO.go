package terminal

type ClientEO struct {
	Id        int64  // 主键
	Name      string // 客户端名称
	LoginIp   string // 客户方ip
	LoginName string // 登录名
	LoginPwd  string // 登录密码
	LoginPort int    // 端口
}
