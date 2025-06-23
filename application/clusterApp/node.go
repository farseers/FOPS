// @area /cluster/
package clusterApp

import (
	"fops/domain/clusterNode"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
)

// NodeList 集群节点列表
// @get nodeList
// @filter application.Jwt
func NodeList() collections.List[docker.DockerNodeVO] {
	lst := clusterNode.NodeList
	lst.Foreach(func(item *docker.DockerNodeVO) {
		// 10s未更新，标记为不健康
		item.IsHealth = time.Since(item.UpdateAt).Seconds() <= 10
	})
	return lst
}
