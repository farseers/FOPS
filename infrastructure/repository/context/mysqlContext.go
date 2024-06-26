package context

import (
	"fops/domain/accountLogin"
	"fops/domain/apps"
	"fops/domain/cluster"
	"fops/domain/configure"
	"fops/infrastructure/repository/model"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/core"
)

// MysqlContext 初始化数据库上下文
var MysqlContext *mysqlContext

// mysqlContext 数据库上下文
type mysqlContext struct {
	// 手动使用事务时必须定义
	core.ITransaction
	// 获取原生ORM框架（不使用TableSet或DomainSet）
	data.IInternalContext
	// 应用
	Apps data.DomainSet[model.AppsPO, apps.DomainObject] `data:"name=apps;migrate;"`
	// cluster表
	Cluster data.DomainSet[model.ClusterPO, cluster.DomainObject] `data:"name=cluster;migrate;"`
	// 节点信息
	ClusterNode data.DomainSet[model.ClusterNodePO, apps.DockerNodeVO] `data:"name=cluster_node;migrate;"`
	// build表
	Build data.DomainSet[model.BuildPO, apps.BuildEO] `data:"name=build;migrate;"`
	// Git
	Git data.DomainSet[model.GitPO, apps.GitEO] `data:"name=git;migrate;"`
	// 配置中心
	Configure data.DomainSet[model.ConfigurePO, configure.DomainObject] `data:"name=configure;migrate;"`
	// 登录帐号
	Login data.DomainSet[model.AccountLoginPO, accountLogin.DomainObject] `data:"name=account_login;migrate;"`
}

// InitMysqlContext 初始化上下文
func InitMysqlContext() {
	MysqlContext = data.NewContext[mysqlContext]("default")
}
