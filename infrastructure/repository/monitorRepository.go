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
)

type monitorRepository struct {
}

// ToListRuleByAppId 获取规则数据
func (receiver *monitorRepository) ToListRuleByAppId(appId string) collections.List[monitor.RuleEO] {
	poList := context.MysqlContext.MonitorRule.Where("app_id = ? and enable = 1", appId).ToList()
	return mapper.ToList[monitor.RuleEO](poList)
}

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
