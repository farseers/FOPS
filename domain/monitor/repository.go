package monitor

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
)

// Repository 仓储接口
type Repository interface {
	// IRepository 通用的仓储接口
	data.IRepository[LogEO]

	// ToListRuleByAppId 查询读取规则
	ToListRuleByAppId(appId string) collections.List[RuleEO]
}
