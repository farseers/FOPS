// @area /terminal/ws/
package terminalApp

import (
	"fops/domain/terminal"

	"github.com/farseer-go/webapi/websocket"
)

// WsSsh 远程命令终端
// @ws ssh
// @filter application.Jwt
func WsSsh(context *websocket.Context[terminal.SshRequest], terminalRepository terminal.Repository) {
	var sshClient terminal.SSHClient
	term := terminal.Terminal{
		Columns: uint32(150),
		Rows:    uint32(30),
	}
	// p为用户输入
	req := context.Receiver()
	// 初始化客户端
	if req.LoginIp != "" {
		info := terminalRepository.ToEntity(req.LoginIp)
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
	sshClient.Connect(context)
}

// WsSshByLogin 远程命令终端
// @ws sshByLogin
// @filter application.Jwt
func WsSshByLogin(context *websocket.Context[terminal.SshRequest], terminalRepository terminal.Repository) {
	var sshClient terminal.SSHClient
	term := terminal.Terminal{
		Columns: uint32(150),
		Rows:    uint32(30),
	}
	// p为用户输入
	req := context.Receiver()
	// 初始化客户端
	if req.IsNotNil() {
		sshClient = terminal.DecodedMsgToSSHClient(req.LoginIp, req.LoginName, req.LoginPwd, req.LoginPort)
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
	sshClient.Connect(context)
}
