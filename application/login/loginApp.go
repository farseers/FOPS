package login

import (
	"fops/application/login/request"
	"fops/application/login/response"
	"fops/domain"
	"fops/domain/accountLogin"
	"time"

	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi"
	"github.com/farseer-go/webapi/check"
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
		"LoginName":  login.LoginName,
		"ClusterIds": login.ClusterIds.ToString(","),
		"DateTime":   time.Now(),
	}
	domain.SetLoginAccount(claims) // 登陆事件 用到
	token, _ := httpContext.Jwt.Build(claims)
	return response.LoginResponse{LoginName: login.LoginName, Token: token}
}

// 修改密码
// @post /user/passport/changePwd
// @filter application.Jwt
func ChangePwd(req request.ChangePwdRequest, accountLoginRepository accountLogin.Repository) {
	check.IsTrue(len(req.LoginPwd) < 6, 403, "密码不能小于6位")
	check.IsTrue(req.LoginPwd != req.ConfirmPwd, 403, "密码不一致")
	curAccount := domain.GetLoginAccount()
	login := accountLoginRepository.ToEntityByAccountName(curAccount.LoginName)
	// 改变密码
	login.ChangeNewPwd(req.LoginPwd)
	err := accountLoginRepository.UpdatePwdByAccountName(login.LoginName, login.LoginPwd)
	exception.ThrowRefuseExceptionError(err)
}
