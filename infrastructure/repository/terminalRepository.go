package repository

import (
	"fops/domain/terminal"
	"github.com/farseer-go/data"
)

type terminalRepository struct {
	// IRepository 通用的仓储接口
	data.IRepository[terminal.ClientEO]
}
