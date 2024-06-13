// @area /cluster/
package clusterApp

import (
	"fops/domain/apps"
	"github.com/farseer-go/collections"
)

// NodeList 集群节点列表
// @get nodeList
// filter application.Jwt
func NodeList(appsRepository apps.Repository) collections.List[apps.DockerNodeVO] {
	lst := appsRepository.GetClusterNodeList()
	lst.Foreach(func(item *apps.DockerNodeVO) {
		item.IsHealth = item.Status == "Ready" && item.Availability == "Active"
	})
	return lst
}
