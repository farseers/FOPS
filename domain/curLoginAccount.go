package domain

import (
	"fops/domain/accountLogin"
	"sync"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/asyncLocal"
	"github.com/farseer-go/fs/parse"
)

// 本地活动时间
var LocalOnline = make(map[int]*sync.Map)

// 本地IP解析
var LocalIPAddress = make(map[string]string)

// 当前登陆账号
var curLoginAccount = asyncLocal.New[accountLogin.DomainObject]()

// 获取当前登陆账号
func GetLoginAccount() accountLogin.DomainObject {
	return curLoginAccount.Get()
}

func SetLoginAccount(claims map[string]any) {
	loginAccount := accountLogin.DomainObject{
		LoginName:  parse.ToString(claims["LoginName"]),
		CreateAt:   parse.ToTime(claims["CreateAt"]),
		ClusterIds: collections.ToList[int](parse.ToString(claims["ClusterIds"])),
	}

	// 权限
	curLoginAccount.Set(loginAccount)
}
