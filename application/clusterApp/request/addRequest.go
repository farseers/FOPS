package request

type AddRequest struct {
	Name           string // 集群名称
	Ip             string // 集群地址
	FScheduleAddr  string // 调度中心地址
	DockerName     string // 仓库名称
	DockerHub      string // 托管地址
	DockerUserName string // 账户名称
	DockerUserPwd  string // 账户密码
	DockerNetwork  string // Docker网络
}
