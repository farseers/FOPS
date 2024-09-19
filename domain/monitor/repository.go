package monitor

import (
	"github.com/farseer-go/collections"
	"time"
)

// Repository 仓储接口
type Repository interface {
	// 规则
	DropDownListAppInfo() collections.List[RuleEO]
	// ToListRuleByAppId 查询读取规则
	ToListRuleByAppId(appId string) collections.List[RuleEO]
	ToListRule() collections.List[RuleEO]
	ToListPageRule(pageSize, pageIndex int) collections.PageList[RuleEO]
	DeleteRule(id int64) error
	ToEntityRule(id int64) RuleEO
	UpdateRule(id int64, rule RuleEO) error
	AddRule(rule RuleEO) error

	// 通知人数据
	// ToListNoticeById 通知人集合
	ToListNoticeById(ids []int) collections.List[NoticeEO]
	ToListPageNotice(pageSize, pageIndex int) collections.PageList[NoticeEO]
	DeleteNotice(id int64) error
	ToEntityNotice(id int64) NoticeEO
	UpdateNotice(id int64, rule NoticeEO) error
	AddRNotice(rule NoticeEO) error

	// 上传数据
	// Save 批量添加监控数据
	Save(lstEO collections.List[DataEO]) error
	// ToListDataByAppIdKey 监控数据
	ToListDataByAppIdKey(appId, key string, top int) collections.List[DataEO]
	// GetMaxTimeByAppId 获取app最大时间
	GetMaxTimeByAppId(appId string) time.Time
	ToListPageData(appId string, pageSize, pageIndex int) collections.PageList[DataEO]

	// 日志
	// SaveLog 批量添加日志数据
	SaveLog(lstEO collections.List[NoticeLogEO]) error
	ToListPageNoticeLog(appId string, pageSize, pageIndex int) collections.PageList[NoticeLogEO]
	DeleteNoticeLog(startTime, endTime time.Time) error
}
