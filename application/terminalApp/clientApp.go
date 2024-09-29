// @area /terminal/
package terminalApp

import (
	"fops/application/terminalApp/response"
	"fops/domain/terminal"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/webapi/check"
)

// ClientList 客户端列表
// @post clientList
// @filter application.Jwt
func ClientList(pageSize, pageIndex int, terminalRepository terminal.Repository) collections.PageList[response.ClientResponse] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	page := terminalRepository.ToPageList(pageSize, pageIndex)
	return mapper.ToPageList[response.ClientResponse](page)
}

// ClientAdd 客户端添加
// @post clientAdd
// @filter application.Jwt
func ClientAdd(req terminal.ClientEO, terminalRepository terminal.Repository) int64 {
	err := terminalRepository.Add(req)
	exception.ThrowWebExceptionError(403, err)
	return req.Id
}

// ClientUpdate 客户端更新
// @post clientUpdate
// @filter application.Jwt
func ClientUpdate(req terminal.ClientEO, terminalRepository terminal.Repository) {
	check.IsTrue(req.Id == 0, 403, "修改数据id参数不能为空")
	info := terminalRepository.ToEntity(req.Id)
	if len(req.LoginPwd) > 0 {
		info.LoginPwd = req.LoginPwd
	}
	info.LoginIp = req.LoginIp
	info.LoginName = req.LoginName
	info.LoginPort = req.LoginPort
	_, err := terminalRepository.Update(req.Id, info)
	exception.ThrowWebExceptionError(403, err)
}

// ClientDel 客户端删除
// @post clientDel
// @filter application.Jwt
func ClientDel(id int64, terminalRepository terminal.Repository) {
	check.IsTrue(id == 0, 403, "删除数据id参数不能为空")
	_, err := terminalRepository.Delete(id)
	exception.ThrowWebExceptionError(403, err)
}

// ClientInfo 客户端详情
// @post clientInfo
// @filter application.Jwt
func ClientInfo(id int64, terminalRepository terminal.Repository) response.ClientResponse {
	check.IsTrue(id == 0, 403, "删除数据id参数不能为空")
	info := terminalRepository.ToEntity(id)
	return mapper.Single[response.ClientResponse](info)
}
