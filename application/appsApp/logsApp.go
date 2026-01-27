// @area /apps/
package appsApp

import (
	"fmt"
	"fops/application/appsApp/response"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/flog"
)

// DockerSwarm 获取容器日志
// @post logs/dockerSwarm
// @filter application.Jwt
func DockerSwarm(appName string, tailCount int) collections.List[response.DockerSwarmResponse] {
	rsp := collections.NewList[response.DockerSwarmResponse]()

	client := docker.NewClient()
	lst := client.Service.PS(appName)
	lst.Foreach(func(item *docker.ServiceTaskVO) {
		// 通过容器id获取日志
		logs, _ := client.Service.Logs(item.ServiceTaskId, tailCount)
		if item.Error != "" {
			containerInspectJson, _ := client.Container.InspectByServiceId(item.ServiceTaskId)
			if len(containerInspectJson) > 0 {
				item.Error = containerInspectJson[0].Status.Err
			}
		}
		serviceLog := logs.First()
		// 没有取到日志时
		if serviceLog.Logs.Count() < 2 {
			serviceLog.Logs.Clear()
			item.Tasks.Foreach(func(taskInstanceVO *docker.TaskInstanceVO) {
				serviceLog.Logs.Add(fmt.Sprintf("%s\t%s\t%s\t%s\t%s", taskInstanceVO.TaskId, taskInstanceVO.Image, taskInstanceVO.Node, taskInstanceVO.State, taskInstanceVO.Error))
			})
		}
		rsp.Add(response.DockerSwarmResponse{
			ServiceTaskVO: *item,
			Log:           flog.ClearColor(serviceLog.Logs.ToString("\r\n")),
		})
	})

	// 所有日志都没有取到时
	if rsp.All(func(item response.DockerSwarmResponse) bool {
		return item.Log == ""
	}) {
		var image string
		inspect, _ := client.Service.Inspect(appName)
		if len(inspect) > 0 {
			image = inspect[0].Spec.TaskTemplate.ContainerSpec.Image
		}

		// 这里取到的是服务日志，即所有容器的日志。需要把他们区分开来
		logs, _ := client.Service.Logs(appName, tailCount*2)
		logs.Foreach(func(serviceLogVO *docker.ServiceLogVO) {
			if curRsp := rsp.Find(func(item *response.DockerSwarmResponse) bool {
				return serviceLogVO.ContainerId == item.ServiceTaskId
			}); curRsp == nil {
				rsp.Add(response.DockerSwarmResponse{
					ServiceTaskVO: docker.ServiceTaskVO{
						ServiceTaskId: serviceLogVO.ContainerId,
						Name:          serviceLogVO.ServiceName,
						Image:         image,
						Node:          serviceLogVO.NodeName,
						State:         "",
						StateInfo:     "",
						Error:         "",
					},
					Log: flog.ClearColor(serviceLogVO.Logs.ToString("\r\n")),
				})
			} else {
				curRsp.Log += flog.ClearColor("\r\n" + serviceLogVO.Logs.ToString("\r\n"))
			}
		})
	}

	// 移除没有日志的项
	rsp.RemoveAll(func(item response.DockerSwarmResponse) bool {
		return item.Log == ""
	})
	return rsp
}
