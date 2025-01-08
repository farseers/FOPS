package repository

import (
	"fops/domain/accountLogin"
	"fops/infrastructure/repository/context"

	"github.com/farseer-go/data"
	"github.com/farseer-go/mapper"
)

type accountLoginRepository struct {
	// IRepository 通用的仓储接口
	data.IRepository[accountLogin.DomainObject]
}

// 根据帐号获取数据
func (receiver *accountLoginRepository) ToEntityByAccountName(loginName string) accountLogin.DomainObject {
	po := context.MysqlContext.Login.Where("login_name = ?", loginName).ToEntity()
	return mapper.Single[accountLogin.DomainObject](po)
}

// 修改密码
func (receiver *accountLoginRepository) UpdatePwdByAccountName(loginName string, loginPwd string) error {
	_, err := context.MysqlContext.Login.Where("login_name = ?", loginName).UpdateValue("login_pwd", loginPwd)
	return err
}
