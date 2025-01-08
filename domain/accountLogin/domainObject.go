package accountLogin

import (
	"fmt"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/utils/encrypt"
)

type DomainObject struct {
	Id         uint64                // 主键
	LoginName  string                // 登陆名称
	LoginPwd   string                // 登录密码
	LoginSalt  string                // 盐
	CreateAt   time.Time             // 创建时间
	ClusterIds collections.List[int] // 集群ID
}

// NewLoginEO 添加登录信息
func NewLoginEO(loginName, loginPwd string) DomainObject {
	salt := parse.RandString(6)
	return DomainObject{
		LoginName: loginName,
		LoginSalt: salt,
		LoginPwd:  encrypt.Md5(salt + loginPwd + salt),
		CreateAt:  time.Now(),
	}
}

func (receiver *DomainObject) IsNil() bool {
	return receiver.LoginName == ""
}

func (receiver *DomainObject) CheckLogin(pwd string) error {
	if receiver.IsNil() {
		return fmt.Errorf("用户名或密码不正确！")
	}
	// 检查密码是否正确
	if receiver.LoginPwd != encrypt.Md5(receiver.LoginSalt+pwd+receiver.LoginSalt) {
		return fmt.Errorf("用户名或密码不正确！")
	}
	return nil
}

// 修改新密码
func (receiver *DomainObject) ChangeNewPwd(pwd string) {
	// 修改新密码
	receiver.LoginPwd = encrypt.Md5(receiver.LoginSalt + pwd + receiver.LoginSalt)
}
