// @area /terminal/
package terminalApp

import (
	"fops/domain/terminal"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/webapi/check"
)

// ClientList 客户端列表
// @post clientList
// @filter application.Jwt
func ClientList(pageSize, pageIndex int, terminalRepository terminal.Repository) collections.PageList[terminal.ClientEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return terminalRepository.ToPageList(pageSize, pageIndex)
}

// ClientAdd 客户端添加
// @post clientAdd
// @filter application.Jwt
func ClientAdd(req terminal.ClientEO, terminalRepository terminal.Repository) {
	err := terminalRepository.Add(req)
	exception.ThrowWebExceptionError(403, err)
}

// ClientUpdate 客户端更新
// @post clientUpdate
// @filter application.Jwt
func ClientUpdate(req terminal.ClientEO, terminalRepository terminal.Repository) {
	check.IsTrue(req.Id == 0, 403, "修改数据id参数不能为空")
	_, err := terminalRepository.Update(req.Id, req)
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
func ClientInfo(id int64, terminalRepository terminal.Repository) terminal.ClientEO {
	check.IsTrue(id == 0, 403, "删除数据id参数不能为空")
	return terminalRepository.ToEntity(id)
}
