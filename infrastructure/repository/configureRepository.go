package repository

import (
	_ "embed"
	"fops/domain/configure"
	"fops/infrastructure/repository/context"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/mapper"
)

type configureRepository struct {
	// IRepository 通用的仓储接口
	data.IRepository[configure.DomainObject]
}

func (receiver *configureRepository) ToListByAppName(appName string) collections.List[configure.DomainObject] {
	sql := "SELECT cp.app_name as app_name, cp.`key` as `key`, cp.value as value, cp.ver as ver FROM configure cp INNER JOIN (SELECT app_name, `key`, MAX(`ver`) AS ver FROM configure where app_name = ? or app_name = 'global' GROUP BY app_name, `key`) AS max_cp ON cp.app_name = max_cp.app_name AND cp.`key` = max_cp.`key` AND cp.Ver = max_cp.ver;"
	lst := context.MysqlContext.Configure.ExecuteSqlToList(sql, appName).ToList()
	return mapper.ToList[configure.DomainObject](lst)
}

func (receiver *configureRepository) ToList() collections.List[configure.DomainObject] {
	sql := "SELECT cp.app_name as app_name, cp.`key` as `key`, cp.value as value, cp.ver as ver FROM configure cp INNER JOIN (SELECT app_name, `key`, MAX(`ver`) AS ver FROM configure GROUP BY app_name, `key`) AS max_cp ON cp.app_name = max_cp.app_name AND cp.`key` = max_cp.`key` AND cp.Ver = max_cp.ver;"
	lst := context.MysqlContext.Configure.ExecuteSqlToList(sql).ToList()
	return mapper.ToList[configure.DomainObject](lst)
}

func (receiver *configureRepository) ToEntity(appName any) configure.DomainObject {
	po := context.MysqlContext.Configure.Where("app_name = ?", appName).Desc("ver").ToEntity()
	return mapper.Single[configure.DomainObject](po)
}

func (receiver *configureRepository) ToEntityByKey(appName, key string) configure.DomainObject {
	po := context.MysqlContext.Configure.Where("app_name = ? and `key` = ?", appName, key).Desc("ver").ToEntity()
	return mapper.Single[configure.DomainObject](po)
}

func (receiver *configureRepository) GetLastVer(appName, key string) int {
	return context.MysqlContext.Configure.Where("app_name = ? and `key` = ?", appName, key).Desc("ver").GetInt("ver")
}

func (receiver *configureRepository) Rollback(appName, key string, ver int) (int64, error) {
	return context.MysqlContext.Configure.Where("app_name = ? and `key` = ? and ver = ?", appName, key, ver).Delete()
}

func (receiver *configureRepository) DeleteKey(appName, key string) (int64, error) {
	return context.MysqlContext.Configure.Where("app_name = ? and `key` = ?", appName, key).Delete()
}
