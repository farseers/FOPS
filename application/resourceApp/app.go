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

			// 总内存、内存使用
			memoryTotal := fmt.Sprintf("%.1fGB", parse.ToFloat64(req.Host.MemoryTotal)/1024/1024/1024)
			memoryUsage := parse.ToFloat64(req.Host.MemoryUsage) / 1024 / 1024
			memoryUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", memoryUsage), 64)

			// 总磁盘、磁盘使用
			diskTotal := fmt.Sprintf("%.1fGB", parse.ToFloat64(req.Host.DiskTotal)/1024/1024/1024)
			diskUsage := parse.ToFloat64(req.Host.DiskUsage) / 1024 / 1024 / 1024
			diskUsage, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", diskUsage), 64)

			// 更新集群节点资源信息
			node := clusterNode.NodeList.Find(func(item *docker.DockerNodeVO) bool {
				return item.IP == req.Host.IP
			})
			if node == nil {
				clusterNode.NodeList.Add(docker.DockerNodeVO{
					IP:            req.Host.IP,
					NodeName:      req.Host.HostName,
					Status:        "Ready",
					Availability:  "Active",
					IsMaster:      req.IsDockerMaster,
					IsHealth:      true,
					EngineVersion: req.DockerEngineVersion,
					OS:            req.Host.OS,
					Architecture:  req.Host.Architecture,
					CPUs:          strconv.Itoa(req.Host.CpuCores),
					Memory:        memoryTotal,
					Label:         collections.List[docker.DockerLabelVO]{},
					UpdateAt:      time.Now(),
				})
				node = clusterNode.NodeList.Find(func(item *docker.DockerNodeVO) bool {
					return item.IP == req.Host.IP
				})
			}

			// 更新集群节点资源信息
			if node != nil {
				node.EngineVersion = req.DockerEngineVersion
				node.OS = req.Host.OS
				node.Architecture = req.Host.Architecture
				node.NodeName = req.Host.HostName
				node.CPUs = strconv.Itoa(req.Host.CpuCores)
				node.Memory = memoryTotal
				node.CpuUsagePercent = req.Host.CpuUsagePercent
				node.MemoryUsagePercent = req.Host.MemoryUsagePercent
				node.MemoryUsage = memoryUsage
				node.Disk = diskTotal
				node.DiskUsagePercent = req.Host.DiskUsagePercent
				node.DiskUsage = diskUsage
				node.UpdateAt = time.Now()
				node.IsHealth = true
			}
		}

		// 更新docker应用资源信息
		if req.Host.IP != "" && req.Dockers.Count() > 0 {
			apps.NodeDockerStatsList.Add(req.Host.IP, req.Dockers)
		}
	})
}
