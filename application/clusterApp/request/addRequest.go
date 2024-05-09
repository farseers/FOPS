package request

type AddRequest struct {
	Name           string // 集群名称
	FopsAddr       string // 集群地址
	FScheduleAddr  string // 调度中心地址
	DockerHub      string // 托管地址
	DockerUserName string // 账户名称
	DockerUserPwd  string // 账户密码
	DockerNetwork  string // Docker网络
}
