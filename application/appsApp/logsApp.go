// @area /apps/
package appsApp

import (
	"fops/domain/apps"
	"github.com/farseer-go/collections"
)

// dockerSwarm 获取容器日志
// @get logs/dockerSwarm
// @filter application.Jwt
func dockerSwarm(appName string, tailCount int, appsIDockerSwarmDevice apps.IDockerSwarmDevice) collections.List[string] {
	return appsIDockerSwarmDevice.Logs(appName, tailCount)
}
