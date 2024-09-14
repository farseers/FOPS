package monitor

import (
	"github.com/farseer-go/collections"
)

// Repository 仓储接口
type Repository interface {
	// ToListRuleByAppId 查询读取规则
	ToListRuleByAppId(appId string) collections.List[RuleEO]
	// 批量添加监控数据
	Save(lstEO collections.List[DataEO]) error
}
