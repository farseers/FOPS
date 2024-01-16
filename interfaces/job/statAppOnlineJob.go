package job

import (
	"encoding/json"
	"fops/domain/apps"
	"fops/domain/register"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/tasks"
	"strings"
)

// StatAppOnlineJob 统计应用的在线实例
func StatAppOnlineJob(*tasks.TaskContext) {
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
		traceContext.Ignore()
	}

	registerRepository := container.Resolve[register.Repository]()
	appsRepository := container.Resolve[apps.Repository]()
	lstGroupBy := registerRepository.StatRegisterApp()

	// 开启事务更新
	container.Resolve[core.ITransaction]("default").Transaction(func() {
		lstApps := appsRepository.ToList()
		lstApps.Foreach(func(item *apps.DomainObject) {
			item.AppName = strings.ToLower(item.AppName)
			// 设置之前的数据，用于比较是否有变化
			setBefore, _ := json.Marshal(item.ActiveInstance)
			if lstRegisterDO, exists := lstGroupBy[strings.ToLower(item.AppName)]; exists {
				item.ActiveInstance = mapper.Array[apps.ActiveInstanceEO](lstRegisterDO)
			} else {
				item.ActiveInstance = make([]apps.ActiveInstanceEO, 0)
			}

			// 有变化，才需用更新数据库
			setAfter, _ := json.Marshal(item.ActiveInstance)
			if string(setBefore) != string(setAfter) {
				_, _ = appsRepository.UpdateActiveInstance(item.AppName, item.ActiveInstance)
			}
		})
	})
}
