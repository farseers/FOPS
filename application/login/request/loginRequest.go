package request

type LoginRequest struct {
	// 登录账号
	LoginName string `validate:"required" label:"登录账号"`
	// 登录密码
	LoginPwd string `validate:"required" label:"登录密码"`
}

type ChangePwdRequest struct {
	// 登录密码
	LoginPwd string `validate:"required" label:"登录密码"`
	// 确认密码
	ConfirmPwd string `validate:"required" label:"确认密码"`
}
