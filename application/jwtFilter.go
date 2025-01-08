package application

import (
	"fops/domain"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/webapi/context"
)

// 账号JWT验证
type Jwt struct {
}

func (receiver Jwt) OnActionExecuting(httpContext *context.HttpContext) {
	traceHand := container.Resolve[trace.IManager]().TraceHand("验证jwt")
	if !httpContext.Jwt.Valid() {
		exception.ThrowWebExceptionf(context.InvalidStatusCode, context.InvalidMessage)
	}

	claims := httpContext.Jwt.GetClaims()

	// 获取登陆账号
	loginName := parse.ToString(claims["LoginName"])
	exception.ThrowWebExceptionBool(len(loginName) == 0, 207, "登录状态失效")

	// 当前登陆账户
	domain.SetLoginAccount(claims)

	traceHand.End(nil)
}

func (receiver Jwt) OnActionExecuted(httpContext *context.HttpContext) {
}
