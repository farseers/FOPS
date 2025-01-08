package clusterNode

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
)

// 当前集群的节点列表（不读库，走本地缓存）
var NodeList = collections.NewList[docker.DockerNodeVO]()
