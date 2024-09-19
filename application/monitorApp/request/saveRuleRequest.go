package request

import (
	"fops/domain/enum/ruleTimeType"
)

type SaveRuleRequest struct {
	Id         int64             // 主键
	AppId      string            // 项目ID
	AppName    string            // 项目名称
	TimeType   ruleTimeType.Enum // 规则时间类型 0小时，1天
	StartTime  string            // 开始时间
	EndTime    string            // 结束时间
	Comparison string            // 比较方式 >  =  <
	KeyName    string            // 监控键
	KeyValue   string            // 监控键值
	Remark     string            // 备注
	Enable     bool              // 是否启用
	NoticeIds  []int             // NoticeIds
}
