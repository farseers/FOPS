package job

import (
	"fops/domain/backupData"
	"time"

	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/tasks"
)

// SyncBackupDataJob 同步备份计划
func SyncBackupDataJob(*tasks.TaskContext) {
	backupDataRepository := container.Resolve[backupData.Repository]()
	do := backupDataRepository.ToNextBackupData()
	if do.IsNil() {
		return
	}
	// 超过60秒的不处理
	s := do.NextBackupAt.Sub(dateTime.Now())
	if s.Seconds() > 60 {
		return
	}
	time.Sleep(s)
	// 时间到了，开始执行
	if err := do.Backup(); err != nil {
		flog.Warning(err.Error())
	}
	backupDataRepository.UpdateAt(do.Id, do)
}
