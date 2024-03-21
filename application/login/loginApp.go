package login

import (
	"fops/application/login/request"
	"fops/application/login/response"
	"fops/domain"
	"fops/domain/accountLogin"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi"
	"time"
)

// 登陆
// @post /user/passport/Login
func Login(req request.LoginRequest, accountLoginRepository accountLogin.Repository) response.LoginResponse {
	httpContext := webapi.GetHttpContext()
	login := accountLoginRepository.ToEntityByAccountName(req.LoginName)
	// 验证
	err := login.CheckLogin(req.LoginPwd)
	exception.ThrowWebExceptionError(403, err)
	claims := map[string]any{
		"LoginName": login.LoginName,
		"DateTime":  time.Now(),
	}
	domain.SetLoginAccount(claims) // 登陆事件 用到
	token, _ := httpContext.Jwt.Build(claims)
	return response.LoginResponse{LoginName: login.LoginName, Token: token}
}
