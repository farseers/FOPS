package job

import (
	"fmt"
	"fops/domain/apps"
	"fops/domain/enum/ruleTimeType"
	"fops/domain/monitor"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/tasks"
	"time"
)

// MonitorRealTimeJob 处理提交过来的监控数据
func MonitorRealTimeJob(*tasks.TaskContext) {
	// 规则
	monitorRepository := container.Resolve[monitor.Repository]()
	appsRepository := container.Resolve[apps.Repository]()
	// 应用规则数据
	ruleList := monitorRepository.ToListRule()
	appNameMap := make(map[string][]monitor.RuleEO)
	ruleList.GroupBy(&appNameMap, func(item monitor.RuleEO) any {
		return item.AppName
	})
	// 添加通知日志
	addLogs := collections.NewList[monitor.NoticeLogEO]()
	// 循环map
	for appName, ruleArray := range appNameMap {
		var ruleInfo monitor.RuleEO

		if appName == "fops" {
			// apps 信息
			appList := appsRepository.ToList()
			// cluster_node 节点信息
			nodeList := appsRepository.GetClusterNodeList()
			for _, rule := range ruleArray {
				ruleInfo = rule
				// 时间类型判断
				var send = false
				switch rule.TimeType {
				case ruleTimeType.Hour:
					startTime := parse.ToInt(rule.StartTime.Format("150405"))
					endTime := parse.ToInt(rule.EndTime.Format("150405"))
					nowTime := parse.ToInt(time.Now().Format("150405"))
					if nowTime >= startTime && nowTime <= endTime {
						send = true
					}
				case ruleTimeType.Day:
					if time.Now().After(rule.StartTime) && time.Now().Before(rule.EndTime) {
						send = true
					}
				}
				if !rule.IsNull() && send {
					comparisonMsg := ""
					switch rule.Comparison {
					case ">":
						switch rule.KeyName {
						case "cpu":
							appList.Foreach(func(app *apps.DomainObject) {
								app.DockerInspect.Foreach(func(docker *apps.DockerInspectVO) {
									if docker.CpuUsagePercent > parse.ToFloat64(rule.KeyValue) {
										comparisonMsg += fmt.Sprintf(" %s %s cpu 大于 %s%", time.Now().Format("01-02 15:04:05"), app.AppName, rule.KeyValue)
									}
								})
							})
							nodeList.Foreach(func(node *docker.DockerNodeVO) {
								if node.CpuUsagePercent > parse.ToFloat64(rule.KeyValue) {
									comparisonMsg += fmt.Sprintf(" %s %s cpu 大于 %s%", time.Now().Format("01-02 15:04:05"), node.NodeName, rule.KeyValue)
								}
							})
						case "memory":
							appList.Foreach(func(app *apps.DomainObject) {
								app.DockerInspect.Foreach(func(docker *apps.DockerInspectVO) {
									if docker.MemoryUsagePercent > parse.ToFloat64(rule.KeyValue) {
										comparisonMsg += fmt.Sprintf(" %s %s memory 大于 %s%", time.Now().Format("01-02 15:04:05"), app.AppName, rule.KeyValue)
									}
								})
							})
							nodeList.Foreach(func(node *docker.DockerNodeVO) {
								if node.MemoryUsagePercent > parse.ToFloat64(rule.KeyValue) {
									comparisonMsg += fmt.Sprintf(" %s %s memory 大于 %s%", time.Now().Format("01-02 15:04:05"), node.NodeName, rule.KeyValue)
								}
							})
						case "disk":
							nodeList.Foreach(func(node *docker.DockerNodeVO) {
								if node.DiskUsagePercent > parse.ToFloat64(rule.KeyValue) {
									comparisonMsg += fmt.Sprintf(" %s %s disk 大于 %s%", time.Now().Format("01-02 15:04:05"), node.NodeName, rule.KeyValue)
								}
							})
						}
					case "<":
						switch rule.KeyName {
						case "cpu":
							appList.Foreach(func(app *apps.DomainObject) {
								app.DockerInspect.Foreach(func(docker *apps.DockerInspectVO) {
									if docker.CpuUsagePercent < parse.ToFloat64(rule.KeyValue) {
										comparisonMsg += fmt.Sprintf(" %s %s cpu 小于 %s%", time.Now().Format("01-02 15:04:05"), app.AppName, rule.KeyValue)
									}
								})
							})
							nodeList.Foreach(func(node *docker.DockerNodeVO) {
								if node.CpuUsagePercent < parse.ToFloat64(rule.KeyValue) {
									comparisonMsg += fmt.Sprintf(" %s %s cpu 小于 %s%", time.Now().Format("01-02 15:04:05"), node.NodeName, rule.KeyValue)
								}
							})
						case "memory":
							appList.Foreach(func(app *apps.DomainObject) {
								app.DockerInspect.Foreach(func(docker *apps.DockerInspectVO) {
									if docker.MemoryUsagePercent < parse.ToFloat64(rule.KeyValue) {
										comparisonMsg += fmt.Sprintf(" %s %s memory 小于 %s%", time.Now().Format("01-02 15:04:05"), app.AppName, rule.KeyValue)
									}
								})
							})
							nodeList.Foreach(func(node *docker.DockerNodeVO) {
								if node.MemoryUsagePercent < parse.ToFloat64(rule.KeyValue) {
									comparisonMsg += fmt.Sprintf(" %s %s memory 小于 %s%", time.Now().Format("01-02 15:04:05"), node.NodeName, rule.KeyValue)
								}
							})
						case "disk":
							nodeList.Foreach(func(node *docker.DockerNodeVO) {
								if node.DiskUsagePercent < parse.ToFloat64(rule.KeyValue) {
									comparisonMsg += fmt.Sprintf(" %s %s disk 小于 %s%", time.Now().Format("01-02 15:04:05"), node.NodeName, rule.KeyValue)
								}
							})
						}
					case "=":
						switch rule.KeyName {
						case "cpu":
							appList.Foreach(func(app *apps.DomainObject) {
								app.DockerInspect.Foreach(func(docker *apps.DockerInspectVO) {
									if docker.CpuUsagePercent == parse.ToFloat64(rule.KeyValue) {
										comparisonMsg += fmt.Sprintf(" %s %s cpu 等于 %s%", time.Now().Format("01-02 15:04:05"), app.AppName, rule.KeyValue)
									}
								})
							})
							nodeList.Foreach(func(node *docker.DockerNodeVO) {
								if node.CpuUsagePercent == parse.ToFloat64(rule.KeyValue) {
									comparisonMsg += fmt.Sprintf(" %s %s cpu 等于 %s%", time.Now().Format("01-02 15:04:05"), node.NodeName, rule.KeyValue)
								}
							})
						case "memory":
							appList.Foreach(func(app *apps.DomainObject) {
								app.DockerInspect.Foreach(func(docker *apps.DockerInspectVO) {
									if docker.MemoryUsagePercent == parse.ToFloat64(rule.KeyValue) {
										comparisonMsg += fmt.Sprintf(" %s %s memory 等于 %s%", time.Now().Format("01-02 15:04:05"), app.AppName, rule.KeyValue)
									}
								})
							})
							nodeList.Foreach(func(node *docker.DockerNodeVO) {
								if node.MemoryUsagePercent == parse.ToFloat64(rule.KeyValue) {
									comparisonMsg += fmt.Sprintf(" %s %s memory 等于 %s%", time.Now().Format("01-02 15:04:05"), node.NodeName, rule.KeyValue)
								}
							})
						case "disk":
							nodeList.Foreach(func(node *docker.DockerNodeVO) {
								if node.DiskUsagePercent == parse.ToFloat64(rule.KeyValue) {
									comparisonMsg += fmt.Sprintf(" %s %s disk 等于 %s%", time.Now().Format("01-02 15:04:05"), node.NodeName, rule.KeyValue)
								}
							})
						}
					}

					// 发送消息 whatsapp
					if len(comparisonMsg) > 0 && len(rule.NoticeIds) > 0 {
						// 通知数据
						noticeList := monitorRepository.ToListNoticeById(rule.NoticeIds)
						noticeList.Foreach(func(not *monitor.NoticeEO) {
							not.Notice(comparisonMsg)
							// 记录日志
							addLogs.Add(monitor.NewLog(rule.AppName, not.Id, not.NoticeType, comparisonMsg))
						})
					}
				}
			}
		} else {

			for _, rule := range ruleArray {
				ruleInfo = rule
				// 获取app数据
				appData := monitorRepository.ToListDataByAppNameKey(appName, rule.KeyName, 1)
				dataInfo := appData.First()
				reqVal := dataInfo.Value
				ruleInfo = rule
				// 时间类型判断
				var send = false
				switch rule.TimeType {
				case ruleTimeType.Hour:
					startTime := parse.ToInt(rule.StartTime.Format("150405"))
					endTime := parse.ToInt(rule.EndTime.Format("150405"))
					nowTime := parse.ToInt(time.Now().Format("150405"))
					if nowTime >= startTime && nowTime <= endTime {
						send = true
					}
				case ruleTimeType.Day:
					if time.Now().After(rule.StartTime) && time.Now().Before(rule.EndTime) {
						send = true
					}
				}
				if !rule.IsNull() && send {
					comparisonMsg := ""
					switch rule.Comparison {
					case ">":
						if parse.ToFloat32(rule.KeyValue) > parse.ToFloat32(reqVal) {
							comparisonMsg = fmt.Sprintf("%s %s 大于 %f", time.Now().Format("01-02 15:04:05"), rule.AppName, parse.ToFloat32(reqVal))
						}
					case "<":
						if parse.ToFloat32(rule.KeyValue) < parse.ToFloat32(reqVal) {
							comparisonMsg = fmt.Sprintf("%s %s 小于 %f", time.Now().Format("01-02 15:04:05"), rule.AppName, parse.ToFloat32(reqVal))
						}
					case "=":
						if parse.ToFloat32(rule.KeyValue) == parse.ToFloat32(reqVal) {
							comparisonMsg = fmt.Sprintf("%s %s 等于 %f", time.Now().Format("01-02 15:04:05"), rule.AppName, parse.ToFloat32(reqVal))
							//comparisonMsg = fmt.Sprintf("%s %s 等于", time.Now().Format("01-02 15:04:05"), rule.AppName)
						}
					}
					// 发送消息 whatsapp
					if len(comparisonMsg) > 0 && len(rule.NoticeIds) > 0 {
						// 通知数据
						noticeList := monitorRepository.ToListNoticeById(rule.NoticeIds)
						noticeList.Foreach(func(not *monitor.NoticeEO) {
							not.Notice(comparisonMsg)
							// 记录日志
							addLogs.Add(monitor.NewLog(rule.AppName, not.Id, not.NoticeType, comparisonMsg))
						})
					}
				}
			}

		}

		// appid 取最大的时间
		maxTime := monitorRepository.GetMaxTimeByAppName(appName)
		// 监控程序是否正常
		if time.Now().Sub(maxTime).Minutes() > 10 {
			time.Sleep(1 * time.Second)
			var comparisonMsg = fmt.Sprintf("%s %s %s", time.Now().Format("01-02 15:04:05"), ruleInfo.AppName, "请检查项目是否已经停止")
			// 通知数据
			noticeList := monitorRepository.ToListNoticeById(ruleInfo.NoticeIds)
			noticeList.Foreach(func(not *monitor.NoticeEO) {
				not.Notice(comparisonMsg)
				// 记录日志
				addLogs.Add(monitor.NewLog(ruleInfo.AppName, not.Id, not.NoticeType, comparisonMsg))
			})
		}
	}

	// 保存日志
	if addLogs.Count() > 0 {
		err := monitorRepository.SaveLog(addLogs)
		exception.ThrowWebExceptionError(403, err)
	}
}
