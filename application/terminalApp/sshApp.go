// @area /terminal/ws/
package terminalApp

import (
	"fops/domain/terminal"
	"github.com/farseer-go/webapi/websocket"
)

// WsSsh 远程命令终端
// @ws ssh
func WsSsh(context *websocket.Context[terminal.SshRequest], terminalRepository terminal.Repository) {
	var sshClient terminal.SSHClient
	term := terminal.Terminal{
		Columns: uint32(150),
		Rows:    uint32(30),
	}
	// p为用户输入
	req := context.Receiver()
	// 初始化客户端
	if req.Id > 0 {
		info := terminalRepository.ToEntity(req.Id)
		if info.Id > 0 {
			sshClient = terminal.DecodedMsgToSSHClient(info.LoginIp, info.LoginName, info.LoginPwd, info.LoginPort)
			err := sshClient.GenerateClient()
			if err != nil {
				sendMap := make(map[string]interface{})
				sendMap["StatusCode"] = "403"
				sendMap["StatusMessage"] = "连接失败"
				_ = context.Send(sendMap)
				return
			} else {
				sendMap := make(map[string]interface{})
				sendMap["StatusCode"] = "200"
				sendMap["StatusMessage"] = "连接成功"
				_ = context.Send(sendMap)
			}
			sshClient.RequestTerminal(term)
		}
	}
	sshClient.Connect(context)
}
