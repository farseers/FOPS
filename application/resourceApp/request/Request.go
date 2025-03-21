package request

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/utils/system"
)

type Request struct {
	Host    system.Resource
	Dockers collections.List[docker.DockerStatsVO]
}
