package monitor

import (
	"github.com/farseer-go/fs/dateTime"
)

type DataEO struct {
	AppId    string            // 项目ID
	AppName  string            // 项目名称
	Key      string            // 监控key
	Value    string            // 监控value
	CreateAt dateTime.DateTime // 发生时间
}

// NewDataEO 新建实体
func NewDataEO(appId, appName string, key, value string) DataEO {
	return DataEO{
		AppId:    appId,
		AppName:  appName,
		Key:      key,
		Value:    value,
		CreateAt: dateTime.Now(),
	}
}
