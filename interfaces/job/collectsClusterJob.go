package job

import (
	"fmt"
	"fops/domain/apps"
	"fops/domain/cluster"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/tasks"
	"github.com/farseer-go/utils/system"
	"github.com/farseer-go/utils/ws"
	"strconv"
	"strings"
	"time"
)

// CollectsClusterJob 3秒收集一次Docker集群信息
func CollectsClusterJob(*tasks.TaskContext) {
	appsRepository := container.Resolve[apps.Repository]()
	clusterRepository := container.Resolve[cluster.Repository]()

	// 收集所有节点的信息
	client := docker.NewClient()
	nodeList := client.Node.List()
	nodeList.Foreach(func(node *docker.DockerNodeVO) {
		vo := client.Node.Info(node.NodeName)
		node.IP = vo.IP
		node.OS = vo.OS
		node.Architecture = vo.Architecture
		node.CPUs = vo.CPUs
		node.Memory = vo.Memory
		node.Label = vo.Label
	})

	// 获取本地集群信息
	localCluster := clusterRepository.GetLocalCluster()
	// 收集所有服务的运行情况
	serviceList := client.Service.List()
	// 如果服务不存在，则添加到列表中，用于更新到数据库中，指明服务的实例为0
	lstApp := appsRepository.ToList()

	// 先把fops中的应用缺少的给补上
	serviceList.Foreach(func(item *docker.ServiceListVO) {
		appDO := lstApp.Find(func(appDO *apps.DomainObject) bool {
			return appDO.AppName == item.Name
		})
		// 本地应用不存在，则添加到fops
		if appDO == nil {
			_ = appsRepository.Add(apps.DomainObject{
				AppName:         item.Name,
				DockerImage:     item.Image,
				DockerInstances: item.Instances,
				DockerReplicas:  item.Replicas,
				IsSys:           true,
			})
		}
	})

	// 遍历所有应用，获取docker swarm副本、实例数量、镜像
	lstApp.Foreach(func(appDO *apps.DomainObject) {
		dockerService := serviceList.Find(func(item *docker.ServiceListVO) bool {
			return item.Name == appDO.AppName
		})
		// 应用没有启用容器服务，跳过
		if dockerService == nil {
			appDO.DockerInstances = 0
			// 系统应用，同时在服务列表中又没有，则删除
			if appDO.IsSys {
				_, _ = appsRepository.Delete(appDO.AppName)
			}
			return
		}

		// 如果是本地集群，则更新镜像信息
		if !localCluster.IsNil() {
			appDO.InitCluster(localCluster.Id)
			appDO.ClusterVer[localCluster.Id].DockerImage = dockerService.Image
		}
		// 当系统应用 或 global模式，才要更新副本数量
		if appDO.IsSys || appDO.DockerNodeRole == "global" {
			appDO.DockerReplicas = dockerService.Replicas
		}
		appDO.DockerInstances = dockerService.Instances
	})

	// 遍历所有应用，得到每个应用的inspect详情
	lstApp.Foreach(func(appDO *apps.DomainObject) {
		// 获取应用的详情
		appDO.DockerInspect = collections.NewList[apps.DockerInspectVO]()
		// 得到该应用正在运行的实例列表
		servicePS := client.Service.PS(appDO.AppName).Where(func(item docker.ServicePsVO) bool {
			return item.State != "Shutdown"
		}).ToList()

		// 遍历每个实例，得到容器ID、IP
		servicePS.Foreach(func(item *docker.ServicePsVO) {
			containerInspectJson, _ := client.Container.InspectByServiceId(item.ServiceId)
			if len(containerInspectJson) == 0 {
				return
			}
			// 匹配对应的节点
			node := nodeList.Find(func(node *docker.DockerNodeVO) bool {
				return node.NodeName == item.Node
			})
			if node == nil {
				node = &docker.DockerNodeVO{}
			}

			// 通过代理节点同步到的容器资源信息
			dockerInspectVO := apps.DockerInspectVO{
				DockerStatsVO: apps.GetDockerStats(containerInspectJson[0].Status.ContainerStatus.ContainerID),
				ServiceID:     item.ServiceId,
				Node:          item.Node,
				NodeIP:        node.IP,
				CreatedAt:     containerInspectJson[0].CreatedAt.Format(time.DateTime),
				UpdatedAt:     containerInspectJson[0].UpdatedAt.Format(time.DateTime),
				State:         containerInspectJson[0].Status.State,
			}

			// IP
			if len(containerInspectJson[0].NetworksAttachments) > 0 && len(containerInspectJson[0].NetworksAttachments[0].Addresses) > 0 {
				dockerInspectVO.ContainerIP = strings.Split(containerInspectJson[0].NetworksAttachments[0].Addresses[0], "/")[0]
			}
			appDO.DockerInspect.Add(dockerInspectVO)

			// 如果应用是fops-agent，则给node节点设置agent的容器IP
			if appDO.AppName == "fops-agent" && dockerInspectVO.ContainerIP != "" {
				node.AgentIP = dockerInspectVO.ContainerIP
				agentNotify <- dockerInspectVO.ContainerIP
			}
		})
	})

	// 通过事务来更新
	container.Resolve[core.ITransaction]("default").Transaction(func() {
		// 更新集群节点信息
		appsRepository.UpdateClusterNode(nodeList)
		// 更新服务运行情况
		if serviceList.Count() > 0 {
			_, _ = appsRepository.UpdateInsReplicas(lstApp)
		}
	})
}

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
				appsRepository.UpdateClusterNodeResourceByAgentIP(agentIP, 0, 0, 0, 0, 0, 0)
				return
			}
			flog.Warningf("接收%s 消息失败：%s", url, err.Error())
			return
		}

		resourceResponse.CpuUsagePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", resourceResponse.CpuUsagePercent), 64)
		resourceResponse.MemoryUsagePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", resourceResponse.MemoryUsagePercent), 64)
		resourceResponse.DiskUsagePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", resourceResponse.DiskUsagePercent), 64)

		memoryUsage := parse.ToFloat64(resourceResponse.MemoryUsage) / 1024 / 1024
		memoryUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", memoryUsage), 64)

		diskUsage := parse.ToFloat64(resourceResponse.DiskUsage) / 1024 / 1024 / 1024
		diskUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", diskUsage), 64)

		// 更新集群节点资源信息
		appsRepository.UpdateClusterNodeResourceByAgentIP(agentIP,
			resourceResponse.CpuUsagePercent,
			resourceResponse.MemoryUsagePercent,
			memoryUsage,
			resourceResponse.DiskTotal/1024/1024/1024,
			resourceResponse.DiskUsagePercent,
			diskUsage,
		)
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