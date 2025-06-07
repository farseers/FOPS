package job

import (
	"fops/domain/linkTrace"
	"fops/domain/logData"
	"time"

	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
)

// SyncBackupDataJob 清除链路数据
func ClearLinkTraceDataJob(*tasks.TaskContext) {
	linkTraceRepository := container.Resolve[linkTrace.Repository]()
	logDataRepository := container.Resolve[logData.Repository]()
	days := configure.GetInt("Fops.LinkTrace.SaveDays")
	if days <= 0 {
		return
	}
	days = days * -1

	// 清除链路数据
	linkTraceRepository.Delete(time.Now().AddDate(0, 0, days))
	linkTraceRepository.DeleteSlow(time.Now().AddDate(0, 0, days))
	logDataRepository.DeleteLog(time.Now().AddDate(0, 0, days))
}
