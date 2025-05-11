package job

import (
	"fops/domain/monitor"
	"strings"

	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/queue"
	"github.com/farseer-go/tasks"
)

// WatchDockerEventJob 监听容器事件
func WatchDockerEventJob(*tasks.TaskContext) {
	// 告警规则数据
	monitorRepository := container.Resolve[monitor.Repository]()
	flog.Infof("监听容器事件")

	dockerClient := docker.NewClient()
	eventResults := dockerClient.Event.Watch()
	for eventResult := range eventResults {
		dateEO := &monitor.DataEO{
			AppName:  eventResult.Actor.Attributes.ComDockerSwarmServiceName,
			Key:      "event",
			Value:    eventResult.Action,
			CreateAt: dateTime.NewUnix(int64(eventResult.Time)),
		}

		// 查找当前应用和KEY的规则
		ruleList := monitorRepository.ToListRule()
		curRuleList := ruleList.Where(func(rule monitor.RuleEO) bool {
			return rule.KeyName == dateEO.Key && strings.Contains(strings.ToLower(rule.AppName), strings.ToLower(dateEO.AppName))
		}).ToList()

		// 比较结果
		curRuleList.Foreach(func(rule *monitor.RuleEO) {
			if rule.CompareResult(dateEO.Value) {
				queue.Push("monitor", dateEO)
			}
		})
	}
}
