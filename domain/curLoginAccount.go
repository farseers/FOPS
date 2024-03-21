package domain

import (
	"fops/domain/accountLogin"
	"github.com/farseer-go/fs/asyncLocal"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/parse"
	"sync"
)

// 本地活动时间
var LocalOnline = make(map[int]*sync.Map)

// 本地IP解析
var LocalIPAddress = make(map[string]string)

// 当前登陆账号
var curLoginAccount = asyncLocal.New[LoginAccountEO]()

// DomainObject 在线用户
type LoginAccountEO struct {
	LoginName string // 登陆名称
}

// 得到领域对象
func (receiver LoginAccountEO) GetDO() accountLogin.DomainObject {
	return container.Resolve[accountLogin.Repository]().ToEntityByAccountName(receiver.LoginName)
}

// 获取当前登陆账号
func GetLoginAccount() LoginAccountEO {
	return curLoginAccount.Get()
}

func SetLoginAccount(claims map[string]any) {
	loginAccount := LoginAccountEO{
		LoginName: parse.ToString(claims["LoginName"]),
	}
	curLoginAccount.Set(loginAccount)
}
