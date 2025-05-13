package repository

import (
	"bytes"
	"fmt"
	"fops/domain/apps"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/snc"
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

func (receiver *appsRepository) ToShortList(isAll bool) collections.List[apps.ShortEO] {
	ts := context.MysqlContext.Apps.Omit("framework_gits", "dockerfile_path", "additional_scripts", "is_sys")
	if !isAll {
		ts.Where("is_sys = false")
	}
	lst := ts.ToList()
	return mapper.ToList[apps.ShortEO](lst)
}

func (repository *buildRepository) ToUTList() collections.List[apps.DomainObject] {
	lst := context.MysqlContext.Apps.Select("app_name", "app_git", "ut_workflows_name").Where("is_sys = 0").ToList()
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
func (receiver *appsRepository) UpdateClusterVer(appName string, dicClusterVer collections.Dictionary[int64, apps.ClusterVerVO]) (int64, error) {
	marshal, _ := dicClusterVer.MarshalJSON()
	return context.MysqlContext.Apps.Where("LOWER(app_name) = ?", appName).UpdateValue("cluster_ver", string(marshal))
}

// UpdateInspect 更新从集群中获取到的实例、副本数量
func (receiver *appsRepository) UpdateInspect(lst collections.List[apps.DomainObject]) (int64, error) {
	sql := bytes.Buffer{}
	sql.WriteString("UPDATE apps SET \n")

	// Instances
	sql.WriteString("docker_instances = case\n")
	lst.Foreach(func(item *apps.DomainObject) {
		sql.WriteString(fmt.Sprintf("when app_name = '%s' then %d\n", item.AppName, item.DockerInstances))
	})
	sql.WriteString("else docker_instances\n")
	sql.WriteString("end \n")

	// Replicas（只更新系统应用或全局应用）
	if lstSysGlobal := lst.Where(func(item apps.DomainObject) bool {
		return item.IsSys || item.DockerNodeRole == "global"
	}).ToList(); lstSysGlobal.Any() {
		sql.WriteString(",docker_replicas = case\n")
		lstSysGlobal.Foreach(func(item *apps.DomainObject) {
			sql.WriteString(fmt.Sprintf("when app_name = '%s' then %d\n", item.AppName, item.DockerReplicas))
		})
		sql.WriteString("else docker_replicas\n")
		sql.WriteString("end \n")
	}
	// cluster_ver
	sql.WriteString(",cluster_ver = case\n")
	lst.Foreach(func(item *apps.DomainObject) {
		marshal, _ := snc.Marshal(item.ClusterVer)
		sql.WriteString(fmt.Sprintf("when app_name = '%s' then '%s'\n", item.AppName, string(marshal)))
	})
	sql.WriteString("else cluster_ver\n")
	sql.WriteString("end \n")

	// docker_inspect
	sql.WriteString(",docker_inspect = case\n")
	lst.Foreach(func(item *apps.DomainObject) {
		marshal, _ := snc.Marshal(item.DockerInspect)
		sql.WriteString(fmt.Sprintf("when app_name = '%s' then '%s'\n", item.AppName, string(marshal)))
	})
	sql.WriteString("else docker_inspect\n")
	sql.WriteString("end \n")

	// where
	sql.WriteString("WHERE 1=1;\n")
	return context.MysqlContext.ExecuteSql(sql.String())
}
