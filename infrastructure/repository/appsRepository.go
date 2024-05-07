package repository

import (
	"encoding/json"
	"fops/domain/apps"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/mapper"
)

type appsRepository struct {
	// IRepository 通用的仓储接口
	data.IRepository[apps.DomainObject]
	buildRepository
	gitRepository
}

func (receiver *appsRepository) UpdateApp(do apps.DomainObject) error {
	po := mapper.Single[model.AppsPO](do)
	_, err := context.MysqlContext.Apps.Where("LOWER(app_name) = ?", po.AppName).Omit("app_name", "docker_ver", "docker_image", "active_instance", "cluster_ver").Update(po)
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
	flog.Info(string(marshal))
	return context.MysqlContext.Apps.Where("LOWER(app_name) = ?", appName).UpdateValue("cluster_ver", string(marshal))
}

func (receiver *appsRepository) UpdateActiveInstance(appName string, eo []apps.ActiveInstanceEO) (int64, error) {
	marshal, _ := json.Marshal(eo)
	return context.MysqlContext.Apps.Where("LOWER(app_name) = ?", appName).UpdateValue("active_instance", string(marshal))
}
