package monitor

import (
	"github.com/farseer-go/collections"
	"time"
)

// Repository 仓储接口
type Repository interface {
	// ToListRuleByAppId 查询读取规则
	ToListRuleByAppId(appId string) collections.List[RuleEO]
	ToListRule() collections.List[RuleEO]
	// ToListNoticeById 通知人集合
	ToListNoticeById(ids []int) collections.List[NoticeEO]
	// Save 批量添加监控数据
	Save(lstEO collections.List[DataEO]) error
	// ToListDataByAppIdKey 监控数据
	ToListDataByAppIdKey(appId, key string, top int) collections.List[DataEO]
	// GetMaxTimeByAppId 获取app最大时间
	GetMaxTimeByAppId(appId string) time.Time
}
