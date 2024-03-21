package request

type LoginRequest struct {
	// 登录账号
	LoginName string `validate:"required" label:"登录账号"`
	// 登录密码
	LoginPwd string `validate:"required" label:"登录密码"`
}
