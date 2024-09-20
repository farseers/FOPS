package monitor

import (
	"github.com/farseer-go/fs/dateTime"
)

type DataEO struct {
	AppName  string            // 项目名称
	Key      string            // 监控key
	Value    string            // 监控value
	CreateAt dateTime.DateTime // 发生时间
}

// NewDataEO 新建实体
func NewDataEO(appName string, key, value string) DataEO {
	return DataEO{
		AppName:  appName,
		Key:      key,
		Value:    value,
		CreateAt: dateTime.Now(),
	}
}
