package job

import (
	"fmt"
	"fops/domain/apps"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/utils/system"
	"github.com/farseer-go/utils/ws"
)

var agentNotify = make(chan string, 100)
var mAgent = make(map[string]string)

// ListenerAgentNotify 监听新的代理节点IP
func ListenerAgentNotify() {
	for {
		agentIP := <-agentNotify

		// 获取主机资源
		if _, exists := mAgent["host_"+agentIP]; !exists {
			mAgent["host_"+agentIP] = ""
			go connectAgentByHostResource(agentIP)
		}

		// 获取容器资源
		if _, exists := mAgent["docker_"+agentIP]; !exists {
			mAgent["docker_"+agentIP] = ""
			go connectAgentByDockerResource(agentIP)
		}
	}
}

// 获取主机资源
func connectAgentByHostResource(agentIP string) {
	appsRepository := container.Resolve[apps.Repository]()

	// 访问获取主机资源
	url := fmt.Sprintf("ws://%s:8888/ws/host/resource", agentIP)
	defer func() {
		delete(mAgent, "host_"+agentIP)
		flog.Debugf("代理节点%s，已断开", url)
	}()

	client, err := ws.Connect(url, 8192)
	client.AutoExit = false

	if err != nil {
		flog.Warningf("连接%s 失败：%s", url, err.Error())
		return
	}

	for {
		var resourceResponse system.Resource
		if err = client.Receiver(&resourceResponse); err != nil {
			if client.IsClose() {
				// 更新集群节点资源信息
				appsRepository.UpdateClusterNodeResourceByAgentIP(agentIP, 0, 0, 0)
				return
			}
			flog.Warningf("接收%s 消息失败：%s", url, err.Error())
			return
		}

		// 更新集群节点资源信息
		appsRepository.UpdateClusterNodeResourceByAgentIP(agentIP,
			resourceResponse.CpuUsagePercent,
			resourceResponse.MemoryUsagePercent,
			resourceResponse.MemoryUsage/1024/1024)
	}
}

// 获取Docker资源
func connectAgentByDockerResource(agentIP string) {
	// 访问获取主机资源
	url := fmt.Sprintf("ws://%s:8888/ws/docker/resource", agentIP)
	defer func() {
		delete(mAgent, "docker_"+agentIP)
		flog.Debugf("代理节点%s，已断开", url)
	}()

	client, err := ws.Connect(url, 8192)
	client.AutoExit = false

	if err != nil {
		flog.Warningf("连接%s 失败：%s", url, err.Error())
		return
	}

	for {
		var resourceResponse collections.List[docker.DockerStatsVO]
		if err = client.Receiver(&resourceResponse); err != nil {
			if client.IsClose() {
				return
			}
			flog.Warningf("接收%s 消息失败：%s", url, err.Error())
			return
		}

		if resourceResponse.Count() > 0 {
			apps.NodeDockerStatsList[agentIP] = resourceResponse
		}
	}
}
