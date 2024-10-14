// @area /cluster/
package clusterApp

import (
	"fops/domain/clusterNode"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
)

// NodeList 集群节点列表
// @get nodeList
// @filter application.Jwt
func NodeList(clusterNodeRepository clusterNode.Repository) collections.List[docker.DockerNodeVO] {
	lst := clusterNodeRepository.GetClusterNodeList()
	lst.Foreach(func(item *docker.DockerNodeVO) {
		item.IsHealth = item.Status == "Ready" && item.Availability == "Active"
	})
	return lst
}
