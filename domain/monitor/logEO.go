package monitor

import (
	"github.com/farseer-go/collections"
	"time"
)

type LogEO struct {
	Id       int64                                  // 主键
	AppId    string                                 // 项目ID
	AppName  string                                 // 项目名称
	Keys     collections.Dictionary[string, string] // 监控键值对
	CreateAt time.Time                              // 发生时间
}

// NewLogEO 新建实体
func NewLogEO(appId, appName string, keys collections.Dictionary[string, string]) LogEO {
	return LogEO{
		AppId:    appId,
		AppName:  appName,
		Keys:     keys,
		CreateAt: time.Now(),
	}
}
