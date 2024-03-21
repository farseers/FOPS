package repository

import (
	"fops/domain/accountLogin"
	"fops/domain/apps"
	"fops/domain/cluster"
	"fops/domain/linkTrace"
	"fops/domain/logData"
	"fops/domain/register"
	"github.com/farseer-go/fs/container"
)

// InitRepository 初始化仓储
func InitRepository() {
	// 应用
	container.Register(func() apps.Repository { return &appsRepository{} })
	// 集群
	container.Register(func() cluster.Repository { return &clusterRepository{} })
	// 链路追踪
	container.Register(func() linkTrace.Repository { return &linkTraceRepository{} })
	// 日志
	container.Register(func() logData.Repository { return &logDataRepository{} })
	// 注册应用
	container.Register(func() register.Repository { return &registerRepository{} })
	// 登录帐号
	container.Register(func() accountLogin.Repository { return &accountLoginRepository{} })
}
