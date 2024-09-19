// @area /terminal/ws/
package terminalApp

import (
	"fops/application/terminalApp/request"
	"fops/domain/terminal"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi/websocket"
	"time"
)

// WsSsh 远程命令终端
// @ws ssh
// @filter application.Jwt
func WsSsh(context *websocket.Context[request.SshRequest], terminalRepository terminal.Repository) {

	var sshClient terminal.SSHClient

	term := terminal.Terminal{
		Columns: uint32(150),
		Rows:    uint32(30),
	}
	context.ReceiverFunc(5*time.Second, func(req *request.SshRequest) {
		// 初始化客户端
		if sshClient.Client == nil && req.Id > 0 {
			info := terminalRepository.ToEntity(req.Id)
			if info.Id > 0 {
				sshClient = terminal.DecodedMsgToSSHClient(info.LoginIp, info.LoginName, info.LoginPwd, info.LoginPort)
				err := sshClient.GenerateClient()
				exception.ThrowWebExceptionError(403, err)
				sshClient.RequestTerminal(term)
			}
		}
		if len(req.Command) > 0 {
			// p为用户输入
			_, err := sshClient.Channel.Write([]byte(req.Command))
			exception.ThrowWebExceptionError(403, err)
			sshClient.Connect(context)
			req.Command = ""
		}

	})

}
