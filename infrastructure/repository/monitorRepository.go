package repository

import (
	"fmt"
	"fops/domain/linkTrace"
	"fops/domain/monitor"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
)

type monitorRepository struct {
}

// ToListRuleByAppName 获取规则数据
func (receiver *monitorRepository) ToListRuleByAppName(appName string) collections.List[monitor.RuleEO] {
	ts := context.MysqlContext.MonitorRule.Where("enable = 1")
	if len(appName) > 0 {
		ts.Where("app_name = ?", appName)
	}
	poList := ts.ToList()
	return mapper.ToList[monitor.RuleEO](poList)
}

func (receiver *monitorRepository) DropDownListAppInfo() collections.List[monitor.RuleEO] {
	sql := `select app_name from monitor_rule where enable=1
			group by  app_name
			order by app_name desc;`
	list := context.MysqlContext.MonitorRule.ExecuteSqlToList(sql)
	return mapper.ToList[monitor.RuleEO](list)
}

// ToListRule 获取所有规则数据
func (receiver *monitorRepository) ToListRule() collections.List[monitor.RuleEO] {
	poList := context.MysqlContext.MonitorRule.Where("enable = 1").ToList()
	return mapper.ToList[monitor.RuleEO](poList)
}

func (receiver *monitorRepository) ToListPageRule(appName string, pageSize, pageIndex int) collections.PageList[monitor.RuleEO] {
	ts := context.MysqlContext.MonitorRule.Desc("id")
	if len(appName) > 0 {
		ts.Where("app_name = ?", appName)
	}
	poList := ts.ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[monitor.RuleEO](poList)
}
func (receiver *monitorRepository) UpdateRuleEnable(id int64, enable bool) error {
	_, err := context.MysqlContext.MonitorRule.Where("id = ?", id).UpdateValue("enable", enable)
	return err
}
func (receiver *monitorRepository) DeleteRule(id int64) error {
	_, err := context.MysqlContext.MonitorRule.Where("id = ?", id).Delete()
	return err
}
func (receiver *monitorRepository) ToEntityRule(id int64) monitor.RuleEO {
	entity := context.MysqlContext.MonitorRule.Where("id = ?", id).ToEntity()
	return mapper.Single[monitor.RuleEO](entity)
}
func (receiver *monitorRepository) UpdateRule(id int64, rule monitor.RuleEO) error {
	po := mapper.Single[model.MonitorRulePO](rule)
	_, err := context.MysqlContext.MonitorRule.Where("id = ?", id).Update(po)
	return err
}
func (receiver *monitorRepository) AddRule(rule monitor.RuleEO) error {
	po := mapper.Single[model.MonitorRulePO](rule)
	err := context.MysqlContext.MonitorRule.Insert(&po)
	return err
}

// ToListNoticeById 通知人集合
func (receiver *monitorRepository) ToListNoticeById(ids []int) collections.List[monitor.NoticeEO] {
	poList := context.MysqlContext.MonitorNotice.Where("id in ? and enable = 1", ids).ToList()
	return mapper.ToList[monitor.NoticeEO](poList)
}
func (receiver *monitorRepository) ToListPageNotice(name string, pageSize, pageIndex int) collections.PageList[monitor.NoticeEO] {
	ts := context.MysqlContext.MonitorNotice.Desc("id")
	if len(name) > 0 {
		ts.Where("name like ?", "%"+name+"%")
	}
	poList := ts.ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[monitor.NoticeEO](poList)
}
func (receiver *monitorRepository) DeleteNotice(id int64) error {
	_, err := context.MysqlContext.MonitorNotice.Where("id = ?", id).Delete()
	return err
}
func (receiver *monitorRepository) ToEntityNotice(id int64) monitor.NoticeEO {
	entity := context.MysqlContext.MonitorNotice.Where("id = ?", id).ToEntity()
	return mapper.Single[monitor.NoticeEO](entity)
}
func (receiver *monitorRepository) UpdateNotice(id int64, rule monitor.NoticeEO) error {
	po := mapper.Single[model.MonitorNoticePO](rule)
	_, err := context.MysqlContext.MonitorNotice.Where("id = ?", id).Update(po)
	return err
}
func (receiver *monitorRepository) AddRNotice(rule monitor.NoticeEO) error {
	po := mapper.Single[model.MonitorNoticePO](rule)
	err := context.MysqlContext.MonitorNotice.Insert(&po)
	return err
}

// Save 保存数据
func (receiver *monitorRepository) Save(lstEO collections.List[monitor.DataEO]) error {
	lstPO := mapper.ToList[model.MonitorDataPO](lstEO)

	if linkTrace.Config.Driver == "clickhouse" {
		// 写入上下文
		if _, err := context.CHContext.MonitorData.InsertList(lstPO, 2000); err != nil {
			_ = flog.Errorf("MonitorData写入ch失败,%s", err.Error())
		}
	} else {
		return fmt.Errorf("不支持的链路追踪驱动：%s", linkTrace.Config.Driver)
	}
	return nil
}

// ToListDataByAppId 监控数据
func (receiver *monitorRepository) ToListDataByAppNameKey(appName, key string, top int) collections.List[monitor.DataEO] {
	poList := context.CHContext.MonitorData.Where("app_name = ? and key = ?", appName, key).Desc("create_at").ToPageList(top, 1)
	return mapper.ToList[monitor.DataEO](poList.List)
}

// 获取app最大时间
func (receiver *monitorRepository) GetMaxTimeByAppName(appName string) time.Time {
	sql := `select max(create_at) from monitor_data where app_name='%s';`
	query := fmt.Sprintf(sql, appName)
	var getTime time.Time
	_, _ = context.CHContext.ExecuteSqlToValue(&getTime, query)
	return getTime
}

func (receiver *monitorRepository) ToListPageData(appName string, pageSize, pageIndex int) collections.PageList[monitor.DataEO] {
	ts := context.CHContext.MonitorData.Desc("create_at")
	if len(appName) > 0 {
		ts.Where("app_name = ?", appName)
	}
	list := ts.ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[monitor.DataEO](list)
}

// SaveLog 批量添加日志数据
func (receiver *monitorRepository) SaveLog(lstEO collections.List[monitor.NoticeLogEO]) error {
	lstPO := mapper.ToList[model.MonitorNoticeLogPO](lstEO)
	_, err := context.MysqlContext.MonitorNoticeLog.InsertList(lstPO, 1000)
	return err
}
func (receiver *monitorRepository) ToListPageNoticeLog(appName string, pageSize, pageIndex int) collections.PageList[monitor.NoticeLogEO] {
	ts := context.MysqlContext.MonitorNoticeLog.Desc("notice_at")
	if len(appName) > 0 {
		ts.Where("app_name = ?", appName)
	}
	poList := ts.ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[monitor.NoticeLogEO](poList)
}
func (receiver *monitorRepository) ToListPageNoticeLogNoRead(top int) collections.List[monitor.NoticeLogEO] {
	ts := context.MysqlContext.MonitorNoticeLog.Where("is_read = 0").Desc("notice_at")
	poList := ts.ToPageList(top, 1)
	return mapper.ToList[monitor.NoticeLogEO](poList.List)
}

func (receiver *monitorRepository) UpdateNoticeLogRead(ids []int) error {
	//_, err := context.MysqlContext.MonitorNoticeLog.Where("is_read = 0 and id in ?", ids).UpdateValue("is_read", true)
	_, err := context.MysqlContext.MonitorNoticeLog.Where("is_read = 0").UpdateValue("is_read", true)
	return err
}

func (receiver *monitorRepository) DeleteNoticeLog(startTime time.Time) error {
	_, err := context.MysqlContext.MonitorNoticeLog.Where("notice_at <= ?", startTime).Delete() //notice_at >= ? and
	return err
}

// 同步时间
func (receiver *monitorRepository) SaveSyncAt(eo monitor.SyncAtEO) error {
	po := mapper.Single[model.MonitorSyncAtPO](eo)
	err := context.MysqlContext.MonitorSyncAt.Insert(&po)
	return err
}
func (receiver *monitorRepository) UpdateSyncAt(appName string, syncAt time.Time) error {
	_, err := context.MysqlContext.MonitorSyncAt.Where("app_name = ?", appName).UpdateValue("sync_at", syncAt)
	return err
}
func (receiver *monitorRepository) ToSyncAtEntity(appName string) monitor.SyncAtEO {
	entity := context.MysqlContext.MonitorSyncAt.Where("app_name = ?", appName).ToEntity()
	return mapper.Single[monitor.SyncAtEO](entity)
}
func (receiver *monitorRepository) IsExistSyncAt(appName string) bool {
	return context.MysqlContext.MonitorSyncAt.Where("app_name = ?", appName).IsExists()
}
