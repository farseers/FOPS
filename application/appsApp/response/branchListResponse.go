package response

import (
	"fops/domain/appsBranch"

	"github.com/farseer-go/collections"
)

type BranchListResponse struct {
	AppName string // 应用名称
	List    collections.List[appsBranch.DomainObject]
}
