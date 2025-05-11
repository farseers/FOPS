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

// Docker核心容器事件映射表
var dockerEventMap = map[string]string{
	"create":        "容器创建完成",
	"start":         "容器启动",
	"die":           "容器进程终止",
	"stop":          "容器被停止",
	"kill":          "容器被强制终止",
	"pause":         "容器暂停",
	"unpause":       "容器恢复运行",
	"restart":       "容器重启",
	"rename":        "容器重命名",
	"destroy":       "容器被删除",
	"update":        "容器配置更新",
	"health_status": "健康检查状态变更",
	"attach":        "附加到容器",
	"detach":        "从容器分离",
}

// WatchDockerEventJob 监听容器事件
func WatchDockerEventJob(*tasks.TaskContext) {
	flog.Infof("监听容器事件")
	// 告警规则数据
	monitorRepository := container.Resolve[monitor.Repository]()

	dockerClient := docker.NewClient()
	eventResults := dockerClient.Event.Watch()
	for eventResult := range eventResults {
		// 过滤其它信息
		if eventResult.Type != "container" || eventResult.Actor.Attributes.ComDockerSwarmServiceName == "" {
			continue
		}

		// 转换成中文事件描述
		if cns, exists := dockerEventMap[eventResult.Action]; exists {
			eventResult.Action = cns
		}

		dateEO := &monitor.DataEO{
			AppName:  eventResult.Actor.Attributes.ComDockerSwarmServiceName,
			Key:      "event",
			Value:    eventResult.Actor.Attributes.Name + "，" + eventResult.Action,
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
