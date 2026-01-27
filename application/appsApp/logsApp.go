// @area /apps/
package appsApp

import (
	"fops/application/appsApp/response"
	"strings"

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
	// 移除第一行的服务信息
	lst.RemoveAt(0)
	lstRunning := lst.Where(func(item docker.TaskInstanceVO) bool {
		return item.State != "Shutdown"
	}).ToList()
	rsp := collections.NewList[response.DockerSwarmResponse]()
	// 有运行中的容器，则只需要正常的容器日志
	if lstRunning.Count() > 0 {
		lst = lstRunning
	}

	lst.Foreach(func(item *docker.TaskInstanceVO) {
		// 通过容器id获取日志
		logs, _ := client.Service.Logs(item.TaskId, tailCount)
		// 有错误时，则通过docker inspect r6r8uboagmln 获取错误详情
		if item.Error != "" {
			containerInspectJson, _ := client.Container.InspectByServiceId(item.TaskId)
			if len(containerInspectJson) > 0 {
				item.Error = containerInspectJson[0].Status.Err
			}
		}
		if item.Error != "" || logs.Count() > 1 {
			rsp.Add(response.DockerSwarmResponse{
				TaskInstanceVO: *item,
				Log:            flog.ClearColor(logs.ToString("\r\n")),
			})
		}
	})

	// 没有取到日志时，取全部
	if lstRunning.Count() == 0 {
		var image string
		inspect, _ := client.Service.Inspect(appName)
		if len(inspect) > 0 {
			image = inspect[0].Spec.TaskTemplate.ContainerSpec.Image
		}

		// 这里取到的是服务日志，即所有容器的日志。需要把他们区分开来
		logs, _ := client.Service.Logs(appName, tailCount*2)
		logs.Foreach(func(item *string) {
			logIndex := strings.Index(*item, "|")
			// 分割不成功，则过滤
			if logIndex <= 0 {
				return
			}
			containerId := strings.TrimSpace((*item)[:logIndex])
			content := strings.TrimSpace((*item)[logIndex+1:])
			if curRsp := rsp.Find(func(item *response.DockerSwarmResponse) bool {
				return strings.Contains(containerId, item.TaskId)
			}); curRsp == nil {
				rsp.Add(response.DockerSwarmResponse{
					TaskInstanceVO: docker.TaskInstanceVO{
						TaskId:    containerId,
						Name:      appName,
						Image:     image,
						Node:      strings.Split(containerId, "@")[1],
						State:     "",
						StateInfo: "",
						Error:     "",
					},
					Log: flog.ClearColor(content),
				})
			} else {
				curRsp.Log += flog.ClearColor("\r\n" + content)
			}
		})
	}
	return rsp
}
