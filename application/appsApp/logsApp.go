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
	client := docker.NewClient()
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
		if logs.Count() < 2 {
			return
		}
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

	// 没有取到日志时，取全部
	if lst.Count() == 0 {
		var image string
		inspect, _ := client.Service.Inspect(appName)
		if len(inspect) > 0 {
			image = inspect[0].Spec.TaskTemplate.ContainerSpec.Image
		}
		logs, _ := client.Service.Logs(appName, tailCount*2)
		rsp.Add(response.DockerSwarmResponse{
			ServicePsVO: docker.ServicePsVO{
				ServiceId: appName,
				Name:      appName,
				Image:     image,
				Node:      "all",
				State:     "所有日志",
				StateInfo: "所有日志",
				Error:     "",
			},
			Log: flog.ClearColor(logs.ToString("\r\n")),
		})
	}
	return rsp
}
