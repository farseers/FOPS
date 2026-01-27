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

		// 没有取到日志时
		if logs.Count() < 2 {
			logs.Clear()
		}
		rsp.Add(response.DockerSwarmResponse{
			ServiceTaskVO: *item,
			Log:           flog.ClearColor(logs.ToString("\r\n")),
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
		logs.Foreach(func(item *string) {
			logIndex := strings.Index(*item, "|")
			// 分割不成功，则过滤
			if logIndex <= 0 {
				return
			}
			containerId := strings.TrimSpace((*item)[:logIndex])
			content := strings.TrimSpace((*item)[logIndex+1:])
			if curRsp := rsp.Find(func(item *response.DockerSwarmResponse) bool {
				return strings.Contains(containerId, item.ServiceTaskId)
			}); curRsp == nil {
				rsp.Add(response.DockerSwarmResponse{
					ServiceTaskVO: docker.ServiceTaskVO{
						ServiceTaskId: containerId,
						Name:          appName,
						Image:         image,
						Node:          strings.Split(containerId, "@")[1],
						State:         "",
						StateInfo:     "",
						Error:         "",
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
