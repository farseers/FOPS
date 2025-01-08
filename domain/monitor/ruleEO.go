package monitor

import (
	"fops/domain/enum/ruleTimeType"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"strconv"
	"strings"
	"time"
)

type RuleEO struct {
	Id          int64             // 主键
	AppName     string            // 应用名称
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
	tips = strings.Replace(tips, "{{Key}}", receiver.KeyName, -1)
	tips = strings.Replace(tips, "{{Value}}", receiver.KeyValue, -1)
	tips = strings.Replace(tips, "{{RealValue}}", realValue, -1)
	tips = strings.Replace(tips, "{{Time}}", time.Now().Format("01-02 15:04:05"), -1)
	return tips
}

// CompareResult 比较结果
func (receiver *RuleEO) CompareResult(reqVal string) bool {
	switch receiver.Comparison {
	case ">":
		return parse.ToFloat32(reqVal) > parse.ToFloat32(receiver.KeyValue)
	case "<":
		return parse.ToFloat32(reqVal) < parse.ToFloat32(receiver.KeyValue)
	case "=":
		if isBool(receiver.KeyValue) {
			return parse.ToBool(receiver.KeyValue) == parse.ToBool(reqVal)
		} else if isFloat(receiver.KeyValue) {
			return parse.ToFloat32(receiver.KeyValue) == parse.ToFloat32(reqVal)
		} else {
			// 字符串判断
			return strings.ToLower(receiver.KeyValue) == strings.ToLower(reqVal)
		}
	case "!=":
		// 字符串判断
		return strings.ToLower(receiver.KeyValue) != strings.ToLower(reqVal)
	}
	return false
}
func isBool(s string) bool {
	s = strings.ToLower(s)
	return s == "true" || s == "false"
}
func isFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
