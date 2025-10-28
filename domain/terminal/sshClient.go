package terminal

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
	"unicode/utf8"

	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi/websocket"
	"golang.org/x/crypto/ssh"
)

type ptyRequestMsg struct {
	Term     string
	Columns  uint32
	Rows     uint32
	Width    uint32
	Height   uint32
	Modelist string
}

type Terminal struct {
	Columns uint32 `json:"cols"`
	Rows    uint32 `json:"rows"`
}

type SSHClient struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	IpAddress string `json:"ipaddress"`
	Port      int    `json:"port"`
	Session   *ssh.Session
	Client    *ssh.Client
	Channel   ssh.Channel
}

// 创建新的ssh客户端时, 默认用户名为root, 端口为22
//
//	func NewSSHClient(ip, userName, passWord string, port int) SSHClient {
//		client := SSHClient{}
//		client.IpAddress = ip
//		client.Username = userName
//		client.Password = passWord
//		client.Port = port
//		// 创建客户端
//		err := client.generateClient()
//		flog.ErrorIfExists(err)
//		// 创建通道
//		client.createChannel()
//		// 开启终端
//		client.openShell()
//		return client
//	}
type SshRequest struct {
	LoginIp   string // 登录IP
	LoginName string // 登录帐号
	LoginPwd  string // 登录密码
	LoginPort int    // 登录端口
	Command   string // 命令
}

func (receiver *SshRequest) IsNotNil() bool {
	return len(receiver.LoginIp) > 0 && len(receiver.LoginName) > 0 && receiver.LoginPort > 0
}

func NewSSHClient() SSHClient {
	client := SSHClient{}
	client.Username = "root"
	client.Port = 22
	return client
}

func DecodedMsgToSSHClient(ip, userName, passWord string, port int) SSHClient {
	client := NewSSHClient()
	client.IpAddress = ip
	client.Username = userName
	client.Password = passWord
	client.Port = port
	return client
}

func (receiver *SSHClient) GenerateClient() error {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		config       ssh.Config
		err          error
	)
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(receiver.Password))
	config = ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}
	clientConfig = &ssh.ClientConfig{
		User:    receiver.Username,
		Auth:    auth,
		Timeout: 5 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	addr = fmt.Sprintf("%s:%d", receiver.IpAddress, receiver.Port)
	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return err
	}
	receiver.Client = client
	return nil
}

func (receiver *SSHClient) RequestTerminal(terminal Terminal) *SSHClient {
	session, err := receiver.Client.NewSession()
	if err != nil {
		log.Println(err)
		return nil
	}
	receiver.Session = session
	channel, inRequests, err := receiver.Client.OpenChannel("session", nil)
	if err != nil {
		log.Println(err)
		return nil
	}
	receiver.Channel = channel
	go func() {
		for req := range inRequests {
			if req.WantReply {
				req.Reply(false, nil)
			}
		}
	}()
	modes := ssh.TerminalModes{
		//ssh.ECHO:          1, // 是否需要回显 1是需要 0不需要
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	var modeList []byte
	for k, v := range modes {
		kv := struct {
			Key byte
			Val uint32
		}{k, v}
		modeList = append(modeList, ssh.Marshal(&kv)...)
	}
	modeList = append(modeList, 0)
	req := ptyRequestMsg{
		Term:     "xterm",
		Columns:  terminal.Columns,
		Rows:     terminal.Rows,
		Width:    uint32(terminal.Columns * 8),
		Height:   uint32(terminal.Columns * 8),
		Modelist: string(modeList),
	}
	ok, err := channel.SendRequest("pty-req", true, ssh.Marshal(&req))
	if !ok || err != nil {
		log.Println(err)
		return nil
	}
	ok, err = channel.SendRequest("shell", true, nil)
	if !ok || err != nil {
		log.Println(err)
		return nil
	}
	return receiver
}

// 连接
func (receiver *SSHClient) Connect(ws *websocket.Context[SshRequest]) {
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
		_ = receiver.Client.Close()
		_ = receiver.Session.Close()
	}()

	//第二个协程将远程主机的返回结果返回给用户
	go func() {
		br := bufio.NewReader(receiver.Channel)
		buf := []byte{}
		t := time.NewTimer(time.Microsecond * 100)
		defer t.Stop()
		// 构建一个信道, 一端将数据远程主机的数据写入, 一段读取数据写入ws
		r := make(chan rune)

		// 另起一个协程, 一个死循环不断的读取ssh channel的数据, 并传给r信道直到连接断开
		go func() {
			catch := exception.Try(func() {
				for {
					x, size, err := br.ReadRune()
					if err != nil {
						log.Println(err)
						//err = ws.Send([]byte("\033[31m已经关闭连接!\033[0m"))
						//exception.ThrowWebExceptionError(403, err)
						if err != nil {
							return
						}
						ws.Close()
						return
					}
					if size > 0 {
						r <- x
					}

				}
			})
			// 处理异常
			catch.CatchException(func(exp any) {
				if exp != nil {
					return
				}
			})
		}()

		// 主循环
		for {
			select {
			// 每隔100微秒, 只要buf的长度不为0就将数据写入ws, 并重置时间和buf
			case <-t.C:
				if len(buf) != 0 {
					err := ws.Send(buf)
					buf = []byte{}
					if err != nil {
						log.Println(err)
						return
					}
				}
				t.Reset(time.Microsecond * 100)
			// 前面已经将ssh channel里读取的数据写入创建的通道r, 这里读取数据, 不断增加buf的长度, 在设定的 100 microsecond后由上面判定长度是否返送数据
			case d := <-r:
				if d != utf8.RuneError {
					p := make([]byte, utf8.RuneLen(d))
					utf8.EncodeRune(p, d)
					buf = append(buf, p...)
				} else {
					buf = append(buf, []byte("@")...)
				}
			case <-ws.Ctx.Done():
				return
			}
		}
	}()
	//这里第一个协程获取用户的输入
	for {
		// p为用户输入
		req := ws.Receiver()
		if len(req.Command) > 0 {
			_, _ = receiver.Channel.Write([]byte(req.Command))
			//exception.ThrowWebExceptionError(403, err)
		}
	}
}
