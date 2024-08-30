package job

import (
	"fops/domain/apps"
	"fops/domain/cluster"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/tasks"
	"strings"
	"time"
)

// CollectsClusterJob 3秒收集一次Docker集群信息
func CollectsClusterJob(*tasks.TaskContext) {
	appsRepository := container.Resolve[apps.Repository]()
	clusterRepository := container.Resolve[cluster.Repository]()

	// 收集所有节点的信息
	client, _ := docker.NewClient()
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

	// 遍历数据库中的应用（包括系统应用）
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
		appDO.DockerReplicas = dockerService.Replicas
		appDO.DockerInstances = dockerService.Instances

		appDO.DockerInspect = make([]apps.DockerInspectVO, 0)
		// 获取应用的详情
		servicePS := client.Service.PS(appDO.AppName)
		servicePS = servicePS.Where(func(item docker.ServicePsVO) bool {
			return item.State != "Shutdown"
		}).ToList()
		servicePS.Foreach(func(item *docker.ServicePsVO) {
			containerInspectJson, _ := client.Container.InspectByServiceId(item.ServiceId)
			if len(containerInspectJson) == 0 {
				return
			}
			dockerInspectVO := apps.DockerInspectVO{
				//ID:        containerInspectJson[0].ID,
				ServiceID:   item.ServiceId,
				ContainerID: containerInspectJson[0].Status.ContainerStatus.ContainerID,
				Node:        item.Node,
				CreatedAt:   containerInspectJson[0].CreatedAt.Format(time.DateTime),
				UpdatedAt:   containerInspectJson[0].UpdatedAt.Format(time.DateTime),
				State:       containerInspectJson[0].Status.State,
			}

			// 使用简短的容器ID
			if len(dockerInspectVO.ContainerID) >= 12 {
				dockerInspectVO.ContainerID = dockerInspectVO.ContainerID[:12]
			}

			// IP
			if len(containerInspectJson[0].NetworksAttachments) > 0 && len(containerInspectJson[0].NetworksAttachments[0].Addresses) > 0 {
				dockerInspectVO.IP = strings.Split(containerInspectJson[0].NetworksAttachments[0].Addresses[0], "/")[0]
			}
			appDO.DockerInspect = append(appDO.DockerInspect, dockerInspectVO)
		})
	})

	container.Resolve[core.ITransaction]("default").Transaction(func() {
		// 更新集群节点信息
		appsRepository.UpdateClusterNode(nodeList)
		// 更新服务运行情况
		if serviceList.Count() > 0 {
			_, _ = appsRepository.UpdateInsReplicas(lstApp)
		}
	})
}
