package monitor

import (
	"github.com/farseer-go/collections"
)

// Repository 仓储接口
type Repository interface {
	// ToListRuleByAppId 查询读取规则
	ToListRuleByAppId(appId string) collections.List[RuleEO]
	ToListRule() collections.List[RuleEO]
	// ToListNoticeById 通知人集合
	ToListNoticeById(ids []string) collections.List[NoticeEO]
	// Save 批量添加监控数据
	Save(lstEO collections.List[DataEO]) error
	// ToListDataByAppId 监控数据
	ToListDataByAppIdKey(appId, key string, top int) collections.List[DataEO]
}
