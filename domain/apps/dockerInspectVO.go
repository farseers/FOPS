package apps

type DockerInspectVO struct {
	ServiceID   string // 服务ID（docker service ps xxx 得到）
	ContainerID string // 容器ID
	Node        string // 节点
	IP          string
	CreatedAt   string
	UpdatedAt   string
	State       string
}
