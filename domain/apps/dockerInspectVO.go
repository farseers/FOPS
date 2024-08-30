package apps

type DockerInspectVO struct {
	ID        string
	Node      string // 节点
	IP        string
	CreatedAt string
	UpdatedAt string
	State     string
}
