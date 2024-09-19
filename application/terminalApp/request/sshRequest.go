package request

type SshRequest struct {
	Id      int64  // 连接ID
	Command string // 输入命令
}
