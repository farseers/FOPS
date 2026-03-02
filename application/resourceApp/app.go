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
func Resource(context *websocket.Context[request.Request]) {
	// 如果appId为空直接返回
	context.ForReceiverFunc(func(req *request.Request) {
		// 更新主机节点资源信息
		if req.Host.CpuUsagePercent > 0 {
			req.Host.MemoryUsagePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", req.Host.MemoryUsagePercent), 64)
			memoryTotal := fmt.Sprintf("%.1fGB", parse.ToFloat64(req.Host.MemoryTotal)/1024/1024/1024)
			memoryUsage, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", parse.ToFloat64(req.Host.MemoryUsage)/1024/1024), 64)

			// 更新集群节点资源信息
			node := clusterNode.NodeList.Find(func(item *docker.DockerNodeVO) bool {
				return item.Status.Addr == req.Host.IP
			})
			if node == nil {
				dockerNodeVO := docker.DockerNodeVO{}
				dockerNodeVO.Status.Addr = req.Host.IP
				dockerNodeVO.Status.State = "Ready"
				dockerNodeVO.Spec.Availability = "Active"
				dockerNodeVO.Spec.Role = "Manager"
				dockerNodeVO.ManagerStatus.Leader = true
				dockerNodeVO.IsHealth = true
				dockerNodeVO.Engine.EngineVersion = req.DockerEngineVersion
				dockerNodeVO.Description.Hostname = req.Host.HostName
				dockerNodeVO.Description.Platform.OS = req.Host.OS
				dockerNodeVO.Description.Platform.Architecture = req.Host.Architecture
				dockerNodeVO.Description.Resources.NanoCPUs = int64(req.Host.CpuCores)
				dockerNodeVO.Description.Resources.MemoryBytes = parse.ToInt64(req.Host.MemoryTotal)
				dockerNodeVO.Description.Resources.Memory = memoryTotal
				dockerNodeVO.Label = collections.NewList[docker.DockerLabelVO]()
				dockerNodeVO.UpdatedAt = time.Now()
				clusterNode.NodeList.Add(dockerNodeVO)

				// 重新排序
				clusterNode.NodeList = clusterNode.NodeList.OrderBy(func(item docker.DockerNodeVO) any {
					return item.Status.Addr
				}).ToList()

				node = clusterNode.NodeList.Find(func(item *docker.DockerNodeVO) bool {
					return item.Status.Addr == req.Host.IP
				})
			}

			// 更新集群节点资源信息
			if node != nil {
				req.Host.CpuUsagePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", req.Host.CpuUsagePercent), 64)

				node.Description.Hostname = req.Host.HostName
				node.Engine.EngineVersion = req.DockerEngineVersion
				node.Description.Platform.OS = req.Host.OS
				node.Description.Platform.Architecture = req.Host.Architecture
				node.Description.Resources.NanoCPUs = int64(req.Host.CpuCores)
				node.Description.Resources.Memory = memoryTotal
				node.Description.Resources.CpuUsagePercent = req.Host.CpuUsagePercent
				node.Description.Resources.MemoryUsagePercent = req.Host.MemoryUsagePercent
				node.Description.Resources.MemoryUsage = memoryUsage
				node.UpdatedAt = time.Now()
				node.IsHealth = true

				var diskList []docker.DiskVO
				var diskTotal uint64
				for _, disk := range req.Host.Disk {
					// 总磁盘、磁盘使用
					diskUsage, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", parse.ToFloat64(disk.DiskUsage)/1024/1024/1024), 64)
					diskUsagePercent, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", disk.DiskUsagePercent), 64)

					diskTotal += disk.DiskTotal
					diskList = append(diskList, docker.DiskVO{
						Path:             disk.Path,
						Disk:             fmt.Sprintf("%.1fGB", parse.ToFloat64(disk.DiskTotal)/1024/1024/1024),
						DiskUsage:        diskUsage,
						DiskUsagePercent: diskUsagePercent,
					})
				}
				node.Description.Resources.DiskTotal = fmt.Sprintf("%.1fGB", parse.ToFloat64(diskTotal)/1024/1024/1024)
				node.Description.Resources.Disk = diskList
			}
		}

		// 更新docker应用资源信息
		if req.Host.IP != "" && req.Dockers.Count() > 0 {
			apps.NodeDockerStatsList.Add(req.Host.IP, req.Dockers)
		}
	})
}
