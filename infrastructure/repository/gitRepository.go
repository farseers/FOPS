package repository

import (
	"fops/domain/apps"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/mapper"
	"time"
)

type gitRepository struct {
}

func (receiver *gitRepository) ToGitEntity(id int64) apps.GitEO {
	po := context.MysqlContext.Git.Where("id = ?", id).ToEntity()
	return mapper.Single[apps.GitEO](po)
}

func (receiver *gitRepository) ToGitList(lstIds collections.List[int64]) collections.List[apps.GitEO] {
	lst := context.MysqlContext.Git.Where("id in (?)", lstIds.ToArray()).ToList()
	return mapper.ToList[apps.GitEO](lst)
}

func (receiver *gitRepository) ToGitListAll(isApp int) collections.List[apps.GitEO] {
	lst := context.MysqlContext.Git.WhereIf(isApp > -1, "is_app = ?", isApp).Order("is_app desc,name asc").ToList()
	return mapper.ToList[apps.GitEO](lst)
}

func (receiver *gitRepository) AddGit(eo apps.GitEO) error {
	po := mapper.Single[model.GitPO](eo)
	return context.MysqlContext.Git.Insert(&po)
}

func (receiver *gitRepository) UpdateGit(eo apps.GitEO) (int64, error) {
	po := mapper.Single[model.GitPO](eo)
	return context.MysqlContext.Git.Where("id = ?", eo.Id).Omit("pull_at").Update(po)
}

func (receiver *gitRepository) DeleteGit(id int64) (int64, error) {
	return context.MysqlContext.Git.Where("id = ?", id).Delete()
}

func (receiver *gitRepository) ExistsGit(id int64) bool {
	return context.MysqlContext.Git.Where("id = ?", id).IsExists()
}

// UpdateForTime 修改GIT的拉取时间
func (receiver *gitRepository) UpdateForTime(id int, pullAt time.Time) (int64, error) {
	return context.MysqlContext.Git.Where("id = ?", id).UpdateValue("pull_at", pullAt)
}
