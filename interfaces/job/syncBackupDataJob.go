package job

import (
	"fops/domain/backupData"
	"time"

	"github.com/farseer-go/fs/container"
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
	s := time.Until(do.NextBackupAt)
	if s.Seconds() > 60 {
		return
	}
	time.Sleep(s)
	// 时间到了，开始执行
	lstBackupHistoryData := do.Backup()
	if lstBackupHistoryData.Count() > 0 {
		cornSchedule, err := backupData.StandardParser.Parse(do.Cron)
		if err != nil {
			flog.Error("同步备份计划时，do.Cron的值不正确导致错误: %s %v", do.Cron, err)
			return
		}
		// 更新时间字段，并生成下一次执行时间。
		do.LastBackupAt = time.Now()
		do.NextBackupAt = cornSchedule.Next(time.Now())
		backupDataRepository.UpdateAt(do.Id, do)

		// 添加历史记录
		backupDataRepository.AddHistory(lstBackupHistoryData)
	}
}
