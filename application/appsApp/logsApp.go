// @area /apps/
package appsApp

import (
	"fops/domain/apps"
)

// DockerSwarm 获取容器日志
// @post logs/dockerSwarm
// @filter application.Jwt
func DockerSwarm(appName string, tailCount int, appsIDockerSwarmDevice apps.IDockerSwarmDevice) string {
	return appsIDockerSwarmDevice.Logs(appName, tailCount).ToString("\r\n")
}
