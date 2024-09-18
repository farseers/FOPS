package repository

import (
	"fmt"
	"fops/domain/linkTrace"
	"fops/domain/monitor"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
	"time"
)

type monitorRepository struct {
}

// ToListRuleByAppId 获取规则数据
func (receiver *monitorRepository) ToListRuleByAppId(appId string) collections.List[monitor.RuleEO] {
	poList := context.MysqlContext.MonitorRule.Where("app_id = ? and enable = 1", appId).ToList()
	return mapper.ToList[monitor.RuleEO](poList)
}

// ToListRule 获取所有规则数据
func (receiver *monitorRepository) ToListRule() collections.List[monitor.RuleEO] {
	poList := context.MysqlContext.MonitorRule.Where("enable = 1").ToList()
	return mapper.ToList[monitor.RuleEO](poList)
}

// ToListDataByAppId 监控数据
func (receiver *monitorRepository) ToListDataByAppIdKey(appId, key string, top int) collections.List[monitor.DataEO] {
	poList := context.CHContext.MonitorData.Where("app_id = ? and key = ?", appId, key).Desc("create_at").ToPageList(top, 1)
	return mapper.ToList[monitor.DataEO](poList.List)
}

// 获取app最大时间
func (receiver *monitorRepository) GetMaxTimeByAppId(appId string) time.Time {
	sql := `select max(create_at) from monitor_data where app_id='%s';`
	query := fmt.Sprintf(sql, appId)
	var getTime time.Time
	_, _ = context.CHContext.ExecuteSqlToValue(&getTime, query)
	return getTime
}

// ToListNoticeById 通知人集合
func (receiver *monitorRepository) ToListNoticeById(ids []int) collections.List[monitor.NoticeEO] {
	poList := context.MysqlContext.MonitorNotice.Where("id in ? and enable = 1", ids).ToList()
	return mapper.ToList[monitor.NoticeEO](poList)
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
