package request

import (
	"fops/domain/apps"

	"github.com/farseer-go/collections"
)

type BuildAddRequest struct {
	AppName                 string
	WorkflowsName           string
	BranchName              string
	EnableBackDefaultBranch bool
	FrameworkList           collections.List[apps.AppsFrameworkEO] // 依赖的框架源代码
}
