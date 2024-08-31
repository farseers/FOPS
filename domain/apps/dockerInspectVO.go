package apps

import "github.com/farseer-go/docker"

type DockerInspectVO struct {
	docker.DockerStatsVO
	ServiceID string // 服务ID（docker service ps xxx 得到）
	Node      string // 节点
	IP        string // 容器IP
	CreatedAt string
	UpdatedAt string
	State     string
}
