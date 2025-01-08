package request

import (
	"github.com/farseer-go/fs/exception"
	"strings"
)

type UpdateRequest struct {
	Id             int64  // 集群Id
	Name           string // 集群名称
	FopsAddr       string // 集群地址
	FScheduleAddr  string // 调度中心地址
	DockerHub      string // 托管地址
	DockerUserName string // 账户名称
	DockerUserPwd  string // 账户密码
	DockerNetwork  string // Docker网络
	IsLocal        bool   // 本地集群
}

func (receiver *UpdateRequest) Check() {
	if receiver.FScheduleAddr != "" {
		exception.ThrowWebExceptionBool(!strings.HasPrefix(strings.ToLower(receiver.FScheduleAddr), "http"), 403, "调度中心地址必须是http开头")
		receiver.FScheduleAddr = strings.TrimSuffix(receiver.FScheduleAddr, "/")
	}
}
