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
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/webapi/websocket"
)

// WsReceive 监控数据接收,更新到本地内存中，供页面展示使用
// @ws resource
func Resource(context *websocket.Context[request.Request]) {
	// 如果appId为空直接返回
	context.ForReceiverFunc(func(req *request.Request) {
		// 主机IP为空直接返回
		if req.Host.IP == "" {
			return
		}

		// 更新主机节点资源信息
		if req.Host.CpuUsagePercent > 0 {
			req.Host.CpuUsagePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", req.Host.CpuUsagePercent), 64)
			req.Host.MemoryUsagePercent, _ = strconv.ParseFloat(fmt.Sprintf("%.1f", req.Host.MemoryUsagePercent), 64)
			memoryTotal := fmt.Sprintf("%.1fGB", parse.ToFloat64(req.Host.MemoryTotal)/1024/1024/1024)
			memoryUsage, _ := strconv.ParseFloat(fmt.Sprintf("%.1f", parse.ToFloat64(req.Host.MemoryUsage)/1024/1024), 64)

			// 更新集群节点资源信息
			node := clusterNode.NodeList.Find(func(item *docker.DockerNodeVO) bool {
				return item.Status.Addr == req.Host.IP
			})
			if node == nil {
				dockerNodeVO := docker.DockerNodeVO{}
				dockerNodeVO.Label = collections.NewList[docker.DockerLabelVO]()
				dockerNodeVO.Status.Addr = req.Host.IP

				flog.Infof("新增集群节点：%s, 角色：%s", req.Host.IP, req.Role)
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
				node.Engine.EngineVersion = req.DockerEngineVersion
				node.ManagerStatus.Leader = req.IsDockerMaster
				node.Spec.Role = req.Role
				node.Description.Hostname = req.Host.HostName
				node.Description.Platform.OS = req.Host.OS
				node.Description.Platform.Architecture = req.Host.Architecture
				node.Description.Resources.NanoCPUs = int64(req.Host.CpuCores)
				node.Description.Resources.MemoryBytes = parse.ToInt64(req.Host.MemoryTotal)
				node.Description.Resources.Memory = memoryTotal
				node.Description.Resources.CpuUsagePercent = req.Host.CpuUsagePercent
				node.Description.Resources.MemoryUsagePercent = req.Host.MemoryUsagePercent
				node.Description.Resources.MemoryUsage = memoryUsage
				node.UpdatedAt = time.Now()
				node.Label = req.Label
				node.Status.State = "Ready"
				node.Spec.Availability = req.Availability
				node.IsHealth = node.Spec.Availability == "Active"

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
		if req.Dockers.Count() > 0 {
			json, _ := snc.Marshal(req.Dockers)
			flog.Infof("接收docker资源监控数据: %s", json)
			apps.NodeDockerStatsList.Add(req.Host.IP, req.Dockers)
		}
	})
}
