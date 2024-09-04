package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"fops/domain/apps"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
)

type appsRepository struct {
	// IRepository 通用的仓储接口
	data.IRepository[apps.DomainObject]
	buildRepository
	gitRepository
}

func (receiver *appsRepository) ToListBySys(isSys bool) collections.List[apps.DomainObject] {
	ts := context.MysqlContext.Apps.Omit("framework_gits", "dockerfile_path", "additional_scripts", "is_sys")
	// 只显示手动添加的应用（不含系统应用）
	if !isSys {
		ts.Where("is_sys = 0")
	}
	lst := ts.ToList()
	return mapper.ToList[apps.DomainObject](lst)
}

func (receiver *appsRepository) UpdateApp(do apps.DomainObject) error {
	po := mapper.Single[model.AppsPO](do)
	_, err := context.MysqlContext.Apps.Where("LOWER(app_name) = ?", po.AppName).Omit("app_name", "docker_ver", "docker_image", "docker_instances", "is_sys").Update(po)
	return err
}

// UpdateDockerVer 修改镜像版本
func (receiver *appsRepository) UpdateDockerVer(appName string, dockerVer int, imageName string) (int64, error) {
	return context.MysqlContext.Apps.Where("LOWER(app_name) = ?", appName).Select("docker_ver", "docker_image").Update(model.AppsPO{
		DockerVer:   dockerVer,
		DockerImage: imageName,
	})
}

// UpdateClusterVer 修改集群的镜像版本
func (receiver *appsRepository) UpdateClusterVer(appName string, dicClusterVer map[int64]*apps.ClusterVerVO) (int64, error) {
	marshal, _ := json.Marshal(dicClusterVer)
	return context.MysqlContext.Apps.Where("LOWER(app_name) = ?", appName).UpdateValue("cluster_ver", string(marshal))
}

// UpdateInsReplicas 更新从集群中获取到的实例、副本数量
func (receiver *appsRepository) UpdateInsReplicas(lst collections.List[apps.DomainObject]) (int64, error) {
	sql := bytes.Buffer{}
	sql.WriteString("UPDATE apps SET \n")

	// Instances
	sql.WriteString("docker_instances = case\n")
	lst.Foreach(func(item *apps.DomainObject) {
		sql.WriteString(fmt.Sprintf("when app_name = '%s' then %d\n", item.AppName, item.DockerInstances))
	})
	sql.WriteString("else docker_instances\n")
	sql.WriteString("end \n")

	// Replicas
	sql.WriteString(",docker_replicas = case\n")
	lst.Foreach(func(item *apps.DomainObject) {
		sql.WriteString(fmt.Sprintf("when app_name = '%s' then %d\n", item.AppName, item.DockerReplicas))
	})
	sql.WriteString("else docker_replicas\n")
	sql.WriteString("end \n")

	// cluster_ver
	sql.WriteString(",cluster_ver = case\n")
	lst.Foreach(func(item *apps.DomainObject) {
		marshal, _ := json.Marshal(item.ClusterVer)
		sql.WriteString(fmt.Sprintf("when app_name = '%s' then '%s'\n", item.AppName, string(marshal)))
	})
	sql.WriteString("else cluster_ver\n")
	sql.WriteString("end \n")

	// docker_inspect
	sql.WriteString(",docker_inspect = case\n")
	lst.Foreach(func(item *apps.DomainObject) {
		marshal, _ := json.Marshal(item.DockerInspect)
		sql.WriteString(fmt.Sprintf("when app_name = '%s' then '%s'\n", item.AppName, string(marshal)))
	})
	sql.WriteString("else docker_inspect\n")
	sql.WriteString("end \n")

	// where
	sql.WriteString("WHERE 1=1;\n")
	return context.MysqlContext.ExecuteSql(sql.String())
}

// UpdateClusterNode 更新集群节点信息
func (receiver *appsRepository) UpdateClusterNode(lst collections.List[docker.DockerNodeVO]) {
	lstPO := mapper.ToList[model.ClusterNodePO](lst)
	lstPO.Foreach(func(item *model.ClusterNodePO) {
		item.UpdateAt = dateTime.Now()
		// 更新数据
		count, err := context.MysqlContext.ClusterNode.Where("node_name", item.NodeName).Omit("cpu_usage_percent", "memory_usage_percent", "memory_usage").Update(*item)
		flog.ErrorIfExists(err)

		// 没有更新到数据时，则插入
		if count == 0 {
			err = context.MysqlContext.ClusterNode.InsertIgnore(item)
			flog.ErrorIfExists(err)
		}
	})
}

func (receiver *appsRepository) GetClusterNodeList() collections.List[docker.DockerNodeVO] {
	lstPO := context.MysqlContext.ClusterNode.Desc("is_master").ToList()
	return mapper.ToList[docker.DockerNodeVO](lstPO)
}

func (receiver *appsRepository) UpdateClusterNodeResourceByAgentIP(agentIP string, cpuUsagePercent, memoryUsagePercent float64, memoryUsage uint64) {
	_, _ = context.MysqlContext.ClusterNode.Where("agent_iP = ?", agentIP).Select("cpu_usage_percent", "memory_usage_percent", "memory_usage").Update(model.ClusterNodePO{
		CpuUsagePercent:    cpuUsagePercent,
		MemoryUsagePercent: memoryUsagePercent,
		MemoryUsage:        memoryUsage,
	})
}
