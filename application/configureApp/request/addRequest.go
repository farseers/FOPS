package request

import "github.com/farseer-go/fs/exception"

type AddRequest struct {
	AppName string // 应用名称
	Key     string // 配置KEY
	Value   string // 配置VALUE
}

func (receiver *AddRequest) Check() {
	if receiver.AppName == "" {
		receiver.AppName = "global"
	}
	exception.ThrowWebExceptionfBool(receiver.Key == "", 403, "Key不能为空")
}
