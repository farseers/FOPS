package monitor

import (
	"fops/domain/enum/ruleTimeType"
	"github.com/farseer-go/collections"
	"strings"
	"time"
)

type RuleEO struct {
	Id          int64             // 主键
	AppName     string            // 应用名称
	TimeType    ruleTimeType.Enum // 规则时间类型 0小时，1天
	StartTime   time.Time         // 开始时间
	EndTime     time.Time         // 结束时间
	Comparison  string            // 比较方式 >  =  <
	KeyName     string            // 监控键
	KeyValue    string            // 监控键值
	Remark      string            // 备注
	TipTemplate string            // 提示模版
	Enable      bool              // 是否启用
	NoticeIds   []int             // NoticeIds
}

// IsNull 判断是否为空
func (receiver *RuleEO) IsNull() bool {
	return receiver == nil || receiver.Id == 0
}

// ToAppNameList 获取appName集合
func (receiver *RuleEO) ToAppNameList() collections.List[string] {
	array := strings.Split(receiver.AppName, ",")
	return collections.NewList[string](array...)
}
func (receiver *RuleEO) GetTipTemplate(appName, realValue string) string {
	var tips string
	tips = strings.Replace(receiver.TipTemplate, "{{AppName}}", appName, -1)
	tips = strings.Replace(receiver.TipTemplate, "{{Key}}", receiver.KeyName, -1)
	tips = strings.Replace(receiver.TipTemplate, "{{Value}}", receiver.KeyValue, -1)
	tips = strings.Replace(receiver.TipTemplate, "{{RealValue}}", realValue, -1)
	tips = strings.Replace(receiver.TipTemplate, "{{Time}}", time.Now().Format("01-02 15:04:05"), -1)
	return tips
}
