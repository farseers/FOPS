package monitor

import (
	"github.com/farseer-go/fs/dateTime"
	"regexp"
	"strings"
)

type DataEO struct {
	AppName  string            // 项目名称
	Key      string            // 监控key
	Value    string            // 监控value
	CreateAt dateTime.DateTime // 发生时间
}

// NewDataEO 新建实体
func NewDataEO(appName string, key, value string) *DataEO {
	return &DataEO{
		AppName:  FilterElement(appName),
		Key:      FilterElement(key),
		Value:    FilterElement(value),
		CreateAt: dateTime.Now(),
	}
}

func FilterElement(val string) string {
	// 去除字符串中的 '\0' 字符
	val = strings.ReplaceAll(val, "\x00", "")
	val = strings.TrimSpace(val)
	val = strings.ReplaceAll(val, " ", "")
	re := regexp.MustCompile(`\s+`)
	val = re.ReplaceAllString(val, "")
	return val
}
