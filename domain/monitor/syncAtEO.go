package monitor

import "time"

type SyncAtEO struct {
	AppName string    // 应用名称
	SyncAt  time.Time // 同步时间
}

// NewSyncAtEO 创建新对象
func NewSyncAtEO(appName string) SyncAtEO {
	return SyncAtEO{
		AppName: appName,
		SyncAt:  time.Now(),
	}
}
