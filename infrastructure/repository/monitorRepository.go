package repository

import (
	"fops/domain/monitor"
	"fops/infrastructure/repository/context"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/mapper"
)

type monitorRepository struct {
	// IRepository 通用的仓储接口
	data.IRepository[monitor.LogEO]
}

// ToListRuleByAppId 获取规则数据
func (receiver *monitorRepository) ToListRuleByAppId(appId string) collections.List[monitor.RuleEO] {
	poList := context.MysqlContext.MonitorRule.Where("app_id = ? and enable = 1", appId).ToList()
	return mapper.ToList[monitor.RuleEO](poList)
}
