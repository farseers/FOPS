// @area /apps/
package appsApp

import (
	"fops/domain/apps"
	"github.com/farseer-go/fs/flog"
)

// DockerSwarm 获取容器日志
// @post logs/dockerSwarm
// @filter application.Jwt
func DockerSwarm(appName string, tailCount int, appsIDockerSwarmDevice apps.IDockerSwarmDevice) string {
	logs := appsIDockerSwarmDevice.Logs(appName, tailCount).ToString("\r\n")
	logs = flog.ClearColor(logs)
	return logs
}
