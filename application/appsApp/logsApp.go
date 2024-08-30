// @area /apps/
package appsApp

import (
	"fops/application/appsApp/response"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/flog"
)

// DockerSwarm 获取容器日志
// @post logs/dockerSwarm
// @filter application.Jwt
func DockerSwarm(appName string, tailCount int) collections.List[response.DockerSwarmResponse] {
	client, _ := docker.NewClient()
	lst := client.Service.PS(appName)
	lstRunning := lst.Where(func(item docker.ServicePsVO) bool {
		return item.State != "Shutdown"
	}).ToList()

	rsp := collections.NewList[response.DockerSwarmResponse]()
	// 有运行中的容器，则只需要正常的容器日志
	if lstRunning.Count() > 0 {
		lst = lstRunning
	}

	lst.Foreach(func(item *docker.ServicePsVO) {
		// 通过容器id获取日志
		logs, _ := client.Service.Logs(item.ServiceId, tailCount)
		// 有错误时，则通过docker inspect r6r8uboagmln 获取错误详情
		if item.Error != "" {
			containerInspectJson, _ := client.Container.InspectByServiceId(item.ServiceId)
			if len(containerInspectJson) > 0 {
				item.Error = containerInspectJson[0].Status.Err
			}
		}

		rsp.Add(response.DockerSwarmResponse{
			ServicePsVO: *item,
			Log:         flog.ClearColor(logs.ToString("\r\n")),
		})
	})
	return rsp
}
