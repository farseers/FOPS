package repository

import (
	"fops/domain/accountLogin"
	"fops/domain/apps"
	"fops/domain/appsBranch"
	"fops/domain/cluster"
	"fops/domain/clusterNode"
	"fops/domain/configure"
	"fops/domain/linkTrace"
	"fops/domain/logData"
	"fops/domain/monitor"
	"fops/domain/terminal"

	"github.com/farseer-go/fs/container"
)

// InitRepository 初始化仓储
func InitRepository() {
	// 应用
	container.Register(func() apps.Repository { return &appsRepository{} })
	// 集群节点
	container.Register(func() clusterNode.Repository { return &clusterNodeRepository{} })
	// 集群
	container.Register(func() cluster.Repository { return &clusterRepository{} })
	// 链路追踪
	container.Register(func() linkTrace.Repository { return &linkTraceRepository{} })
	// 日志
	container.Register(func() logData.Repository { return &logDataRepository{} })
	// 登录帐号
	container.Register(func() accountLogin.Repository { return &accountLoginRepository{} })
	// 配置中心
	container.Register(func() configure.Repository { return &configureRepository{} })
	// 监控中心
	container.Register(func() monitor.Repository { return &monitorRepository{} })
	// ssh客户端
	container.Register(func() terminal.Repository { return &terminalRepository{} })
	// 应用分支
	container.Register(func() appsBranch.Repository { return &appsBranchRepository{} })
}
