// @area /apps/
package appsApp

import (
	"fops/application/appsApp/response"
	"fops/domain/apps"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
)

// DockerSwarm 获取容器日志
// @post logs/dockerSwarm
// @filter application.Jwt
func DockerSwarm(appName string, tailCount int, appsIDockerSwarmDevice apps.IDockerSwarmDevice) collections.List[response.DockerSwarmResponse] {
	lst := appsIDockerSwarmDevice.PS(appName)
	lstRunning := lst.Where(func(item apps.DockerInstanceVO) bool {
		return item.State == "Running"
	}).ToList()

	rsp := collections.NewList[response.DockerSwarmResponse]()
	// 有运行中的容器，则只需要正常的容器日志
	if lstRunning.Count() > 0 {
		lst = lstRunning
	}

	lst.Foreach(func(item *apps.DockerInstanceVO) {
		// 通过容器id获取日志
		logs := appsIDockerSwarmDevice.Logs(item.Id, tailCount).ToString("\r\n")
		logs = flog.ClearColor(logs)

		rsp.Add(response.DockerSwarmResponse{
			DockerInstanceVO: *item,
			Log:              logs,
		})
	})
	return rsp
}
