package terminal

import "github.com/farseer-go/data"

// Repository 仓储接口
type Repository interface {
	// IRepository 通用的仓储接口
	data.IRepository[ClientEO]
}
