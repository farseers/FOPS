package repository

import (
	"fops/domain/_/eumBuildStatus"
	"fops/domain/_/eumBuildType"
	"fops/domain/apps"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
)

type buildRepository struct {
}

// GetBuildNumber 获取构建的编号
func (repository *buildRepository) GetBuildNumber(appName string) int {
	return context.MysqlContext.Build.Where("LOWER(app_name) = ?", appName).Order("build_number desc").GetInt("build_number")
}

func (repository *buildRepository) AddBuild(eo *apps.BuildEO) error {
	po := mapper.Single[model.BuildPO](eo)
	err := context.MysqlContext.Build.Insert(&po)
	eo.Id = po.Id
	return err
}
func (repository *buildRepository) ToBuildList(appName string, buildType eumBuildType.Enum, pageSize int, pageIndex int) collections.PageList[apps.BuildEO] {
	ts := context.MysqlContext.Build.Desc("id").Where("build_type = ?", buildType)
	// 筛选appName
	appName = strings.TrimSpace(appName)
	if appName != "" {
		ts.Where("LOWER(app_name) = ?", appName)
	}
	lstPO := ts.Select("id", "app_name", "build_number", "status", "is_success", "create_at", "finish_at", "workflows_name", "branch_name").ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[apps.BuildEO](lstPO)
}

// GetUnBuildInfo 获取未构建的任务
func (repository *buildRepository) GetUnBuildInfo(buildType eumBuildType.Enum) apps.BuildEO {
	po := context.MysqlContext.Build.Where("status = ? and build_type = ?", eumBuildStatus.None, buildType).Asc("id").ToEntity() //  and build_server_id = ?	, core.AppId
	return mapper.Single[apps.BuildEO](po)
}

// SetBuilding 设置任务为构建中
func (repository *buildRepository) SetBuilding(id int64) {
	_, _ = context.MysqlContext.Build.Where("id = ?", id).Select("status", "create_at").Update(model.BuildPO{
		Status:   eumBuildStatus.Building,
		CreateAt: time.Now(),
	})
}

// UpdateBuilding 更新构建任务
func (repository *buildRepository) UpdateBuilding(id int64, env apps.EnvVO) {
	_, _ = context.MysqlContext.Build.Where("id = ?", id).Select("cluster_id", "docker_image", "branch_name").Update(model.BuildPO{
		ClusterId:   env.ClusterId,
		DockerImage: env.DockerImage,
		BranchName:  env.BranchName,
	})
}

// SetSuccess 任务完成
func (repository *buildRepository) SetSuccess(id int64, env apps.EnvVO) {
	_, _ = context.MysqlContext.Build.Where("id = ?", id).Select("status", "is_success", "finish_at", "env", "docker_image").Update(model.BuildPO{
		Status:      eumBuildStatus.Finish,
		IsSuccess:   true,
		FinishAt:    time.Now(),
		Env:         env,
		DockerImage: env.DockerImage,
	})
}

// SetSuccessForFops 任务完成
func (repository *buildRepository) SetSuccessForFops(id int64) {
	_, _ = context.MysqlContext.Build.Where("id = ?", id).Select("status", "is_success", "finish_at").Update(model.BuildPO{
		Status:    eumBuildStatus.Finish,
		IsSuccess: true,
		FinishAt:  time.Now(),
	})
}

// SetCancel 主动取消任务
func (repository *buildRepository) SetCancel(id int64, env apps.EnvVO) {
	_, _ = context.MysqlContext.Build.Where("id = ?", id).Select("status", "is_success", "finish_at", "env", "docker_image").Update(model.BuildPO{
		Status:      eumBuildStatus.Cancel,
		IsSuccess:   false,
		FinishAt:    time.Now(),
		Env:         env,
		DockerImage: env.DockerImage,
	})
}

// SetFail 任务失败
func (repository *buildRepository) SetFail(id int64, env apps.EnvVO) {
	_, _ = context.MysqlContext.Build.Where("id = ? and status <> ?", id, eumBuildStatus.Cancel).Select("status", "is_success", "finish_at", "env", "docker_image").Update(model.BuildPO{
		Status:      eumBuildStatus.Finish,
		IsSuccess:   false,
		FinishAt:    time.Now(),
		Env:         env,
		DockerImage: env.DockerImage,
	})
}

// GetStatus 获取构建状态
func (repository *buildRepository) GetStatus(id int64) eumBuildStatus.Enum {
	return eumBuildStatus.Enum(context.MysqlContext.Build.Where("id = ?", id).GetInt("status"))
}

// UpdateFailDockerImage 更新构建中状态的构建记录
func (repository *buildRepository) UpdateFailDockerImage(appName string, dockerImage string) (int64, error) {
	return context.MysqlContext.Build.Select("status", "is_success", "finish_at").Where("app_name = ? and status = ? and docker_image = ?", appName, eumBuildStatus.Building, dockerImage).
		Update(model.BuildPO{
			Status:    eumBuildStatus.Finish,
			IsSuccess: true,
			FinishAt:  time.Now(),
		})
}

func (repository *buildRepository) GetLastBuilding(buildType eumBuildType.Enum) apps.BuildEO {
	po := context.MysqlContext.Build.Where("status = ? and build_type = ?", eumBuildStatus.Building, buildType).Desc("id").ToEntity()
	return mapper.Single[apps.BuildEO](po)
}

func (repository *buildRepository) ToBuildEntity(id int64) apps.BuildEO {
	po := context.MysqlContext.Build.Where("id = ?", id).ToEntity()
	return mapper.Single[apps.BuildEO](po)
}
