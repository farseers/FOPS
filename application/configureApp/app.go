// @area /configure/
package configureApp

import (
	"fops/domain/configure"
	"github.com/farseer-go/collections"
)

// List 获取应用的配置列表
// @post list
func List(appName string, configureRepository configure.Repository) collections.List[configure.DomainObject] {
	return configureRepository.ToListByAppName(appName)
}
