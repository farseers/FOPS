// @area /cluster/
package clusterApp

import (
	"fops/domain/apps"
	"github.com/farseer-go/collections"
)

// NodeList 集群节点列表
// @get nodeList
// @filter application.Jwt
func NodeList(appsRepository apps.Repository) collections.List[apps.DockerNodeVO] {
	return appsRepository.GetClusterNodeList()
}
