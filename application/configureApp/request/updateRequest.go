package request

import "github.com/farseer-go/fs/exception"

type UpdateRequest struct {
	AppName string // 应用名称
	Key     string // 配置KEY
	Value   string // 配置VALUE
}

func (receiver *UpdateRequest) Check() {
	exception.ThrowWebExceptionfBool(receiver.AppName == "", 403, "应用名称不能为空")
	exception.ThrowWebExceptionfBool(receiver.Key == "", 403, "Key不能为空")
}
