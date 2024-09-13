// @area /cluster/
package clusterApp

import (
	"fops/domain/apps"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
)

// NodeList 集群节点列表
// @get nodeList
// @filter application.Jwt
func NodeList(appsRepository apps.Repository) collections.List[docker.DockerNodeVO] {
	lst := appsRepository.GetClusterNodeList()
	lst.Foreach(func(item *docker.DockerNodeVO) {
		item.IsHealth = item.Status == "Ready" && item.Availability == "Active"
	})
	return lst
}
