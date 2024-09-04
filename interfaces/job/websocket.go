package job

import (
	"fmt"
	"fops/domain/apps"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/utils/system"
	"github.com/farseer-go/utils/ws"
)

var agentNotify = make(chan string, 100)
var mAgent = make(map[string]string)

func ListenerAgentNotify() {
	for {
		ip := <-agentNotify

		// 获取主机资源
		if _, exists := mAgent["host_"+ip]; !exists {
			mAgent["host_"+ip] = ""
			go connectAgentByHostResource(ip)
		}

		// 获取容器资源
		//if _, exists := mAgent["docker_"+ip]; !exists {
		//	mAgent["docker_"+ip] = ""
		//	go connectAgentByHostResource(ip)
		//}
	}
}

// 获取主机资源
func connectAgentByHostResource(ip string) {
	defer delete(mAgent, "host_"+ip)

	// 访问获取主机资源
	url := fmt.Sprintf("ws://%s:8888/ws/host/resource", ip)
	client, err := ws.Connect(url, 1024)
	if err != nil {
		flog.Warningf("连接%s 失败：%s", url, err.Error())
		return
	}

	flog.Infof("已连接新的代理节点%s，持续接收消息中", url)
	appsRepository := container.Resolve[apps.Repository]()
	for {
		var resourceResponse core.ApiResponse[system.Resource]
		if err = client.Receiver(&resourceResponse); err != nil {
			flog.Warningf("接收%s 消息失败：%s", url, err.Error())
			return
		}

		// 更新集群节点资源信息
		appsRepository.UpdateClusterNodeResourceByAgentIP(ip,
			resourceResponse.Data.CpuUsagePercent,
			resourceResponse.Data.MemoryUsagePercent,
			resourceResponse.Data.MemoryUsage/1024/1024)
	}
}
