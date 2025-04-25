package job

import (
	"fmt"
	"fops/domain/apps"
	"fops/domain/clusterNode"
	"fops/domain/monitor"
	"runtime/debug"
	"strings"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/queue"
	"github.com/farseer-go/tasks"
)

// MonitorFopsJob 监控fops数据
func MonitorFopsJob(*tasks.TaskContext) {
	appsRepository := container.Resolve[apps.Repository]()
	// 规则
	monitorRepository := container.Resolve[monitor.Repository]()
	// apps 信息
	appList := appsRepository.ToList()
	addMonitorData := collections.NewList[monitor.DataEO]()
	// 应用数据
	appList.Foreach(func(app *apps.DomainObject) {
		app.DockerInspect.Foreach(func(docker *apps.DockerInspectVO) {
			addMonitorData.Add(monitor.DataEO{
				AppName:  app.AppName,
				Key:      "cpu",
				Value:    parse.ToString(docker.CpuUsagePercent),
				CreateAt: dateTime.Now(),
			})
			addMonitorData.Add(monitor.DataEO{
				AppName:  app.AppName,
				Key:      "memory",
				Value:    parse.ToString(docker.MemoryUsagePercent),
				CreateAt: dateTime.Now(),
			})
			addMonitorData.Add(monitor.DataEO{
				AppName:  app.AppName,
				Key:      "instances",
				Value:    parse.ToString(app.DockerInstances == app.DockerReplicas),
				CreateAt: dateTime.Now(),
			})
		})
	})
	// 节点数据
	clusterNode.NodeList.Foreach(func(node *docker.DockerNodeVO) {
		addMonitorData.Add(monitor.DataEO{
			AppName:  fmt.Sprintf("%s(%s)", node.IP, node.NodeName),
			Key:      "cpu",
			Value:    parse.ToString(node.CpuUsagePercent),
			CreateAt: dateTime.Now(),
		})
		addMonitorData.Add(monitor.DataEO{
			AppName:  fmt.Sprintf("%s(%s)", node.IP, node.NodeName),
			Key:      "memory",
			Value:    parse.ToString(node.MemoryUsagePercent),
			CreateAt: dateTime.Now(),
		})
		// 多个硬盘路径，取最大占用值
		diskUsagePercent := collections.NewList(node.Disk...).Max(func(item docker.DiskVO) any {
			return item.DiskUsagePercent
		}).(float64)
		addMonitorData.Add(monitor.DataEO{
			AppName:  fmt.Sprintf("%s(%s)", node.IP, node.NodeName),
			Key:      "disk",
			Value:    parse.ToString(diskUsagePercent),
			CreateAt: dateTime.Now(),
		})
		addMonitorData.Add(monitor.DataEO{
			AppName:  fmt.Sprintf("%s(%s)", node.IP, node.NodeName),
			Key:      "pcStatus",
			Value:    node.Status,
			CreateAt: dateTime.Now(),
		})
		addMonitorData.Add(monitor.DataEO{
			AppName:  fmt.Sprintf("%s(%s)", node.IP, node.NodeName),
			Key:      "nodeAvailability",
			Value:    node.Availability,
			CreateAt: dateTime.Now(),
		})
	})
	// 应用规则数据
	ruleList := monitorRepository.ToListRule()
	appNameList := collections.NewList[string]()
	// 添加消息队列
	addMonitorData.Foreach(func(item *monitor.DataEO) {
		curRuleList := ruleList.Where(func(rule monitor.RuleEO) bool {
			return rule.KeyName == item.Key && strings.Contains(strings.ToLower(rule.AppName), strings.ToLower(item.AppName))
		}).ToList()
		curRuleList.Foreach(func(rule *monitor.RuleEO) {
			if rule.CompareResult(item.Value) {
				queue.Push("monitor", item)
			}
		})
		appNameList.Add(item.AppName)
	})
	// 立即释放内存返回给操作系统
	debug.FreeOSMemory()
}
