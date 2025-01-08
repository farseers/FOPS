package request

import (
	"fops/domain/enum/ruleTimeType"
)

type SaveRuleRequest struct {
	Id          int64             // 主键
	AppId       string            // 项目ID
	AppName     string            // 项目名称
	TimeType    ruleTimeType.Enum // 规则时间类型 0小时，1天
	StartDate   string            // 开始小时
	EndDate     string            // 结束小时
	StartDay    string            // 开始天
	EndDay      string            // 结束天
	Comparison  string            // 比较方式 >  =  <
	KeyName     string            // 监控键
	KeyValue    string            // 监控键值
	Remark      string            // 备注
	TipTemplate string            // 提示模版
	Enable      bool              // 是否启用
	NoticeIds   []int             // NoticeIds
}
