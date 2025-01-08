package monitor

import (
	"time"

	"github.com/farseer-go/collections"
)

// Repository 仓储接口
type Repository interface {
	// 规则
	DropDownListAppInfo() collections.List[RuleEO]
	// ToListRuleByAppId 查询读取规则
	ToListRuleByAppName(appName string) collections.List[RuleEO]
	ToListRule() collections.List[RuleEO]
	ToListPageRule(appName string, pageSize, pageIndex int) collections.PageList[RuleEO]
	DeleteRule(id int64) error
	ToEntityRule(id int64) RuleEO
	UpdateRule(id int64, rule RuleEO) error
	AddRule(rule RuleEO) error

	// 通知人数据
	// ToListNoticeById 通知人集合
	ToListNoticeById(ids []int) collections.List[NoticeEO]
	ToListPageNotice(name string, pageSize, pageIndex int) collections.PageList[NoticeEO]
	DeleteNotice(id int64) error
	ToEntityNotice(id int64) NoticeEO
	UpdateNotice(id int64, rule NoticeEO) error
	AddRNotice(rule NoticeEO) error

	// 上传数据
	// Save 批量添加监控数据
	Save(lstEO collections.List[DataEO]) error
	// ToListDataByAppNameKey 监控数据
	ToListDataByAppNameKey(appName, key string, top int) collections.List[DataEO]
	// GetMaxTimeByAppName 获取app最大时间
	GetMaxTimeByAppName(appName string) time.Time
	ToListPageData(appName string, pageSize, pageIndex int) collections.PageList[DataEO]
	// 日志
	// SaveLog 批量添加日志数据
	SaveLog(lstEO collections.List[NoticeLogEO]) error
	ToListPageNoticeLog(appName string, pageSize, pageIndex int) collections.PageList[NoticeLogEO]
	DeleteNoticeLog(startTime time.Time) error
	ToListPageNoticeLogNoRead(top int) collections.List[NoticeLogEO]
	UpdateNoticeLogRead(ids []int) error
	// 同步时间
	SaveSyncAt(eo SyncAtEO) error
	ToSyncAtEntity(appName string) SyncAtEO
	IsExistSyncAt(appName string) bool
	UpdateSyncAt(appName string, syncAt time.Time) error
}
