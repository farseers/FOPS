package job

import (
	"fops/domain/backupData"
	"time"

	"github.com/farseer-go/fs/container"
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
	s := time.Until(do.NextBackupAt)
	if s.Seconds() > 60 {
		return
	}
	time.Sleep(s)

	// 时间到了，开始执行
	lstBackupHistoryData := do.Backup()
	if lstBackupHistoryData.Count() > 0 {
		// 更新时间字段，并生成下一次执行时间。
		do.LastBackupAt = time.Now()
		cornSchedule, _ := backupData.StandardParser.Parse(do.Cron)
		do.NextBackupAt = cornSchedule.Next(time.Now())
		backupDataRepository.Update(do.Id, do)

		// 添加历史记录
		backupDataRepository.AddHistory(lstBackupHistoryData)
	}
}
