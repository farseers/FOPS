package job

import (
	"fmt"
	"fops/domain/monitor"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/tasks"
)

// MonitorRealTimeJob 处理提交过来的监控数据
func MonitorRealTimeJob(*tasks.TaskContext) {
	// 规则
	monitorRepository := container.Resolve[monitor.Repository]()
	// 应用规则数据
	ruleList := monitorRepository.ToListRule()
	// 添加通知日志
	addLogs := collections.NewList[monitor.NoticeLogEO]()
	// 循环map
	ruleList.Foreach(func(rule *monitor.RuleEO) {
		// 应用名称集合
		rule.ToAppNameList().Foreach(func(appName *string) {
			// appid 取最大的时间
			maxTime := monitorRepository.ToSyncAtEntity(*appName)
			// 监控程序是否正常
			if time.Since(maxTime.SyncAt).Minutes() > 10 {
				time.Sleep(1 * time.Second)
				var comparisonMsg = fmt.Sprintf("%s %s %s", time.Now().Format("01-02 15:04:05"), *appName, "请检查项目是否已经停止")
				// 通知数据
				noticeList := monitorRepository.ToListNoticeById(rule.NoticeIds)
				noticeList.Foreach(func(not *monitor.NoticeEO) {
					not.Notice(comparisonMsg)
					// 记录日志
					addLogs.Add(monitor.NewLog(*appName, not.Id, not.Name, not.NoticeType, comparisonMsg))
				})
			}

		})

	})

	// 保存日志
	if addLogs.Count() > 0 {
		err := monitorRepository.SaveLog(addLogs)
		exception.ThrowWebExceptionError(403, err)
	}
}
