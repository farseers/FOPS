package job

import (
	"fops/domain/apps"
	"fops/domain/monitor"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/tasks"
)

// MonitorFopsJob 监控fops数据
func MonitorFopsJob(*tasks.TaskContext) {
	appsRepository := container.Resolve[apps.Repository]()
	monitorRepository := container.Resolve[monitor.Repository]()
	// apps 信息
	appList := appsRepository.ToList()
	// cluster_node 节点信息
	nodeList := appsRepository.GetClusterNodeList()
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
		})
	})
	// 节点数据
	nodeList.Foreach(func(node *docker.DockerNodeVO) {
		addMonitorData.Add(monitor.DataEO{
			AppName:  node.NodeName,
			Key:      "cpu",
			Value:    parse.ToString(node.CpuUsagePercent),
			CreateAt: dateTime.Now(),
		})
		addMonitorData.Add(monitor.DataEO{
			AppName:  node.NodeName,
			Key:      "memory",
			Value:    parse.ToString(node.MemoryUsagePercent),
			CreateAt: dateTime.Now(),
		})
		addMonitorData.Add(monitor.DataEO{
			AppName:  node.NodeName,
			Key:      "disk",
			Value:    parse.ToString(node.DiskUsagePercent),
			CreateAt: dateTime.Now(),
		})
	})
	err := monitorRepository.Save(addMonitorData)
	exception.ThrowWebExceptionError(403, err)
}