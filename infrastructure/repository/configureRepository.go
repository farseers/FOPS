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
	sql := "SELECT * FROM configure cp INNER JOIN (SELECT app_name, `key`, MAX(`ver`) AS ver FROM configure where app_name = '?' or app_name = 'global' GROUP BY app_name, `key`) AS max_cp ON cp.app_name = max_cp.app_name AND cp.`key` = max_cp.`key` AND cp.Ver = max_cp.ver;"
	lst := context.MysqlContext.Configure.ExecuteSqlToList(sql, appName).ToList()
	return mapper.ToList[configure.DomainObject](lst)
}
