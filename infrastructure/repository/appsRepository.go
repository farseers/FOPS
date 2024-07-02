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

func (receiver *appsRepository) ToList() collections.List[apps.DomainObject] {
	lst := context.MysqlContext.Apps.Omit("framework_gits", "dockerfile_path", "additional_scripts").ToList()
	return mapper.ToList[apps.DomainObject](lst)
}

func (receiver *appsRepository) UpdateApp(do apps.DomainObject) error {
	po := mapper.Single[model.AppsPO](do)
	_, err := context.MysqlContext.Apps.Where("LOWER(app_name) = ?", po.AppName).Omit("app_name", "docker_ver", "docker_image", "docker_instances").Update(po)
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

	// where
	sql.WriteString("WHERE 1=1;\n")
	return context.MysqlContext.ExecuteSql(sql.String())
}

// UpdateClusterNode 更新集群节点信息
func (receiver *appsRepository) UpdateClusterNode(lst collections.List[apps.DockerNodeVO]) {
	lstPO := mapper.ToList[model.ClusterNodePO](lst)
	lstPO.Foreach(func(item *model.ClusterNodePO) {
		item.UpdateAt = dateTime.Now()
		// 更新数据
		count, err := context.MysqlContext.ClusterNode.Where("node_name", item.NodeName).Update(*item)
		flog.ErrorIfExists(err)

		// 没有更新到数据时，则插入
		if count == 0 {
			err = context.MysqlContext.ClusterNode.InsertIgnore(item)
			flog.ErrorIfExists(err)
		}
	})
}

func (receiver *appsRepository) GetClusterNodeList() collections.List[apps.DockerNodeVO] {
	lstPO := context.MysqlContext.ClusterNode.Desc("is_master").ToList()
	return mapper.ToList[apps.DockerNodeVO](lstPO)
}
