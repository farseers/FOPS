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
func NodeList() collections.List[docker.DockerNodeVO] {
	return clusterNode.NodeList
}
