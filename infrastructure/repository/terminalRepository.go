package repository

import (
	"fops/domain/terminal"
	"fops/infrastructure/repository/context"

	"github.com/farseer-go/data"
)

type terminalRepository struct {
	// IRepository 通用的仓储接口
	data.IRepository[terminal.ClientEO]
}

// ToListRuleByAppName 获取规则数据
func (receiver *terminalRepository) ExistsHost(ip string) bool {
	ts := context.MysqlContext.TerminalClient.Where("login_ip = ?", ip)
	return ts.IsExists()
}
