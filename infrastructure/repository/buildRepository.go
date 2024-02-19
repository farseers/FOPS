package repository

import (
	"fops/domain/_/eumBuildStatus"
	"fops/domain/apps"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"strings"
	"time"
)

type buildRepository struct {
}

// GetBuildNumber 获取构建的编号
func (repository *buildRepository) GetBuildNumber(appName string) int {
	return context.MysqlContext.Build.Where("LOWER(app_name) = ?", appName).Order("Id desc").GetInt("build_number")
}

func (repository *buildRepository) AddBuild(eo apps.BuildEO) error {
	po := mapper.Single[model.BuildPO](eo)
	return context.MysqlContext.Build.Insert(&po)
}

func (repository *buildRepository) ToBuildList(appName string, pageSize int, pageIndex int) collections.PageList[apps.BuildEO] {
	ts := context.MysqlContext.Build.Desc("id")
	// 筛选appName
	appName = strings.TrimSpace(appName)
	if appName != "" {
		ts.Where("LOWER(app_name) = ?", appName)
	}
	lstPO := ts.ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[apps.BuildEO](lstPO)
}

// GetUnBuildInfo 获取未构建的任务
func (repository *buildRepository) GetUnBuildInfo() apps.BuildEO {
	po := context.MysqlContext.Build.Where("status = ?", eumBuildStatus.None).Asc("id").ToEntity() //  and build_server_id = ?	, core.AppId
	return mapper.Single[apps.BuildEO](po)
}

// SetBuilding 设置任务为构建中
func (repository *buildRepository) SetBuilding(id int64) {
	_, _ = context.MysqlContext.Build.Where("id = ?", id).Select("status", "create_at").Update(model.BuildPO{
		Status:   eumBuildStatus.Building,
		CreateAt: time.Now(),
	})
}

// SetSuccess 任务完成
func (repository *buildRepository) SetSuccess(id int64) {
	_, _ = context.MysqlContext.Build.Where("id = ?", id).Select("status", "is_success", "finish_at").Update(model.BuildPO{
		Status:    eumBuildStatus.Finish,
		IsSuccess: true,
		FinishAt:  time.Now(),
	})
}

// SetCancel 主动取消任务
func (repository *buildRepository) SetCancel(id int64) {
	_, _ = context.MysqlContext.Build.Where("id = ?", id).Select("status", "is_success", "finish_at").Update(model.BuildPO{
		Status:    eumBuildStatus.Finish,
		IsSuccess: false,
		FinishAt:  time.Now(),
	})
}

// GetStatus 获取构建状态
func (repository *buildRepository) GetStatus(id int64) eumBuildStatus.Enum {
	return eumBuildStatus.Enum(context.MysqlContext.Build.Where("id = ?", id).GetInt("status"))
}
