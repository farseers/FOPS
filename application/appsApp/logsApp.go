// @area /apps/
package appsApp

import (
	"fops/domain/apps"
	"github.com/farseer-go/collections"
)

// DockerSwarm 获取容器日志
// @post logs/dockerSwarm
// @filter application.Jwt
func DockerSwarm(appName string, tailCount int, appsIDockerSwarmDevice apps.IDockerSwarmDevice) collections.List[string] {
	return appsIDockerSwarmDevice.Logs(appName, tailCount)
}
