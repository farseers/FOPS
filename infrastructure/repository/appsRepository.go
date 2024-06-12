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
	"github.com/farseer-go/mapper"
	"strings"
)

type appsRepository struct {
	// IRepository 通用的仓储接口
	data.IRepository[apps.DomainObject]
	buildRepository
	gitRepository
}

func (receiver *appsRepository) UpdateApp(do apps.DomainObject) error {
	po := mapper.Single[model.AppsPO](do)
	_, err := context.MysqlContext.Apps.Where("LOWER(app_name) = ?", po.AppName).Omit("app_name", "docker_ver", "docker_image", "active_instance").Update(po)
	return err
}

// UpdateDockerVer 修改镜像版本
func (receiver *appsRepository) UpdateDockerVer(appName string, dockerVer int, imageName string) (int64, error) {
	_, _ = context.MysqlContext.Apps.Where("LOWER(app_name) = ?", appName).UpdateValue("docker_ver", dockerVer)
	return context.MysqlContext.Apps.Where("LOWER(app_name) = ?", appName).UpdateValue("docker_image", imageName)
}

// UpdateClusterVer 修改集群的镜像版本
func (receiver *appsRepository) UpdateClusterVer(appName string, dicClusterVer map[int64]*apps.ClusterVerVO) (int64, error) {
	marshal, _ := json.Marshal(dicClusterVer)
	return context.MysqlContext.Apps.Where("LOWER(app_name) = ?", appName).UpdateValue("cluster_ver", string(marshal))
}

// UpdateInsReplicas 更新从集群中获取到的实例、副本数量
func (receiver *appsRepository) UpdateInsReplicas(lst collections.List[apps.DockerName]) (int64, error) {
	var appNames []string
	lst.Select(&appNames, func(item apps.DockerName) any {
		return item.Name
	})

	sql := bytes.Buffer{}
	sql.WriteString("UPDATE apps SET \n")

	// Instances
	sql.WriteString("docker_instances = case\n")
	lst.Foreach(func(item *apps.DockerName) {
		sql.WriteString(fmt.Sprintf("when app_name = '%s' then %d\n", item.Name, item.Instances))
	})
	sql.WriteString("else docker_instances\n")
	sql.WriteString("end \n")

	// Replicas
	sql.WriteString(",docker_replicas = case\n")
	lst.Foreach(func(item *apps.DockerName) {
		sql.WriteString(fmt.Sprintf("when app_name = '%s' then %d\n", item.Name, item.Replicas))
	})
	sql.WriteString("else docker_replicas\n")
	sql.WriteString("end \n")

	// where
	sql.WriteString(fmt.Sprintf("WHERE app_name in ('%s');\n", strings.Join(appNames, "','")))
	return context.MysqlContext.ExecuteSql(sql.String())
}
