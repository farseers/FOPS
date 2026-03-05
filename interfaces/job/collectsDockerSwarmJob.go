package job

import (
	"fops/domain/apps"
	"fops/domain/cluster"
	"fops/domain/clusterNode"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/tasks"
)

// CollectsDockerSwarmJob 3秒收集一次Docker集群信息
func CollectsDockerSwarmJob(*tasks.TaskContext) {
	appsRepository := container.Resolve[apps.Repository]()
	clusterRepository := container.Resolve[cluster.Repository]()

	// 收集所有节点的信息
	dockerClient := docker.NewClient()

	// 获取本地集群信息
	localCluster := clusterRepository.GetLocalCluster()
	// 收集所有服务的运行情况
	serviceList := dockerClient.Service.List()
	// 所有节点信息
	lstNode := dockerClient.Node.List()
	// 没有读取到应用，则退出
	if serviceList.Count() == 0 {
		return
	}
	// 如果服务不存在，则添加到列表中，用于更新到数据库中，指明服务的实例为0
	lstApp := appsRepository.ToList()
	// 没有读取到应用，则退出
	if lstApp.Count() == 0 {
		return
	}

	// 先把fops中缺少的应用的给补上
	serviceList.Foreach(func(item *docker.ServiceListVO) {
		appDO := lstApp.Find(func(appDO *apps.DomainObject) bool {
			return appDO.AppName == item.Spec.Name
		})
		// 本地应用不存在，则添加到fops
		if appDO == nil {
			_ = appsRepository.Add(apps.DomainObject{
				AppName:         item.Spec.Name,
				DockerInstances: 0,
				DockerReplicas:  item.Spec.Mode.Replicated.Replicas,
				IsSys:           true,
			})
		}
	})

	// 遍历所有应用，更新实际的docker swarm副本、镜像
	lstApp.Foreach(func(appDO *apps.DomainObject) {
		dockerService := serviceList.Find(func(item *docker.ServiceListVO) bool {
			return item.Spec.Name == appDO.AppName
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
			appDO.ClusterVer.Update(localCluster.Id, func(value *apps.ClusterVerVO) {
				value.DockerImage = dockerService.Spec.TaskTemplate.ContainerSpec.Image
			})
		}
		// 当系统应用 或 global模式，才要更新副本数量
		if appDO.IsSys {
			appDO.DockerReplicas = dockerService.Spec.Mode.Replicated.Replicas
		}
		if appDO.DockerNodeRole == "global" {
			appDO.DockerReplicas = lstNode.Count()
		}
	})

	// 遍历所有应用，得到每个应用的inspect详情
	lstApp.Foreach(func(appDO *apps.DomainObject) {
		// 获取应用的详情
		if appDO.DockerInspect.Count() == 0 {
			appDO.DockerInspect = collections.NewList[apps.DockerInspectVO]()
		}

		// 根据ServiceId获取所有实例
		allInstanceList := dockerClient.Service.PS(lstNode, appDO.AppName)
		if allInstanceList.Count() == 0 {
			return
		}
		// 只保留运行中的实例
		runInstanceList := allInstanceList.Where(func(item docker.ServiceTaskVO) bool {
			return item.State == "running" //return item.State != "Shutdown"
		}).ToList()

		// 清空不存在的实例列表
		appDO.DockerInspect.RemoveAll(func(dockerInspectVO apps.DockerInspectVO) bool {
			return !runInstanceList.Where(func(serviceVO docker.ServiceTaskVO) bool {
				return dockerInspectVO.TaskId == serviceVO.ServiceTaskId
			}).Any()
		})

		// 遍历每个实例，得到容器ID、IP
		runInstanceList.Foreach(func(serviceVO *docker.ServiceTaskVO) {
			// 服务已存在于本地实例列表中，则跳过
			curInstance := appDO.DockerInspect.Find(func(dockerInspectVO *apps.DockerInspectVO) bool {
				return dockerInspectVO.TaskId == serviceVO.ServiceTaskId
			})
			// 匹配对应的节点
			node := clusterNode.NodeList.Find(func(node *docker.DockerNodeVO) bool {
				return node.Description.Hostname == serviceVO.NodeName
			})
			// 实例存在，则只更新资源信息
			if node != nil && curInstance != nil {
				curInstance.DockerStatsVO = apps.GetDockerStats(node.Status.Addr, serviceVO.ServiceTaskId)
				return
			}

			// 只有当节点是健康状态才加入到实例列表中。
			if node != nil && node.IsHealth {
				// 读取单个实例的详情
				containerInspectJson, _ := dockerClient.Task.Inspect(serviceVO.ServiceTaskId)
				if containerInspectJson.ID != "" {
					// 通过代理节点同步到的容器资源信息
					dockerInspectVO := apps.DockerInspectVO{
						DockerStatsVO: apps.GetDockerStats(node.Status.Addr, serviceVO.ServiceTaskId),
						NodeID:        serviceVO.NodeID,
						NodeName:      serviceVO.NodeName,
						NodeIP:        serviceVO.NodeIP,
						CreatedAt:     containerInspectJson.CreatedAt.Format(time.DateTime),
						UpdatedAt:     containerInspectJson.UpdatedAt.Format(time.DateTime),
						State:         containerInspectJson.Status.State,
					}

					// 容器IP
					if len(containerInspectJson.NetworksAttachments) > 0 && len(containerInspectJson.NetworksAttachments[0].Addresses) > 0 {
						dockerInspectVO.ContainerIP = strings.Split(containerInspectJson.NetworksAttachments[0].Addresses[0], "/")[0]
					}
					appDO.DockerInspect.Add(dockerInspectVO)
				}
			}
		})
		// 实例数量
		appDO.DockerInstances = runInstanceList.Count()
		time.Sleep(100 * time.Millisecond)
	})

	// 更新服务运行情况
	_, _ = appsRepository.UpdateInspect(lstApp)
}
