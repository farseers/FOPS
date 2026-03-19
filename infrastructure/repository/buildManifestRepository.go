package repository

import (
	"fops/domain/apps"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
)

type buildManifestRepository struct {
}

// AddBuildManifestBatch 批量添加构建清单
func (repository *buildManifestRepository) AddBuildManifestBatch(lst collections.List[apps.BuildManifestEO]) error {
	if lst.Count() == 0 {
		return nil
	}

	lstPO := mapper.ToList[model.BuildManifestPO](lst)
	_, err := context.MysqlContext.BuildManifest.InsertList(lstPO, 1000)
	return err
}

// GetLastBuilds 获取应用最近N次构建记录（仅应用自身的构建，不含依赖）
func (repository *buildManifestRepository) GetLastBuilds(appName string, limit int) collections.List[apps.BuildManifestEO] {
	lstPO := context.MysqlContext.BuildManifest.
		Where("app_name = ? AND git_name = ?", appName, appName).
		Desc("create_at").
		Limit(limit).
		ToList()
	return mapper.ToList[apps.BuildManifestEO](lstPO)
}

// GetManifestsByDockerImage 根据镜像获取构建清单（包含应用和所有依赖）
func (repository *buildManifestRepository) GetManifestsByDockerImage(dockerImage string) collections.List[apps.BuildManifestEO] {
	lstPO := context.MysqlContext.BuildManifest.
		Where("docker_image = ?", dockerImage).
		Asc("git_name"). // 排序：应用名在前，依赖库按字母排序
		ToList()
	return mapper.ToList[apps.BuildManifestEO](lstPO)
}
