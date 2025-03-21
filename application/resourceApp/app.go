// @area /ws/
package resourceApp

import (
	"fmt"
	"fops/application/resourceApp/request"
	"fops/domain/apps"
	"fops/domain/clusterNode"
	"strconv"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/webapi/websocket"
)

// WsReceive 监控数据接收
// @ws resource
func Resource(context *websocket.Context[request.Request], clusterNodeRepository clusterNode.Repository) {
	// 如果appId为空直接返回
	context.ForReceiverFunc(func(req *request.Request) {
		// 更新主机节点资源信息
		if req.Host.CpuUsagePercent > 0 {
			req.Host.CpuUsagePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", req.Host.CpuUsagePercent), 64)
			req.Host.MemoryUsagePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", req.Host.MemoryUsagePercent), 64)
			req.Host.DiskUsagePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", req.Host.DiskUsagePercent), 64)

			memoryUsage := parse.ToFloat64(req.Host.MemoryUsage) / 1024 / 1024
			memoryUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", memoryUsage), 64)

			diskUsage := parse.ToFloat64(req.Host.DiskUsage) / 1024 / 1024 / 1024
			diskUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", diskUsage), 64)

			// 更新集群节点资源信息
			node := clusterNode.NodeList.Find(func(item *docker.DockerNodeVO) bool {
				return item.IP == req.Host.IP
			})
			if node == nil {
				clusterNode.NodeList.Add(docker.DockerNodeVO{
					IP:            req.Host.IP,
					NodeName:      "",
					Status:        "",
					Availability:  "",
					IsMaster:      false,
					IsHealth:      true,
					EngineVersion: "",
					OS:            "",
					Architecture:  "",
					CPUs:          strconv.Itoa(req.Host.CpuCores),
					Memory:        "",
					Label:         collections.List[docker.DockerLabelVO]{},
				})
				node = clusterNode.NodeList.Find(func(item *docker.DockerNodeVO) bool {
					return item.IP == req.Host.IP
				})
			}

			// 更新集群节点资源信息
			if node != nil {
				node.CpuUsagePercent = req.Host.CpuUsagePercent
				node.MemoryUsagePercent = req.Host.MemoryUsagePercent
				node.MemoryUsage = memoryUsage
				node.Disk = req.Host.DiskTotal / 1024 / 1024 / 1024
				node.DiskUsagePercent = req.Host.DiskUsagePercent
				node.DiskUsage = diskUsage
			}
		}

		// 更新docker应用资源信息
		if req.Host.IP != "" && req.Dockers.Count() > 0 {
			apps.NodeDockerStatsList[req.Host.IP] = req.Dockers
		}

		// 间隔更新
		if time.Now().Second()%30 == 0 {
			// 更新集群节点信息
			clusterNodeRepository.UpdateClusterNode(clusterNode.NodeList)
		}
	})
}
