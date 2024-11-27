package response

import (
	"fops/domain/_/eumBuildStatus"
)

type BuildListResponse struct {
	Id            int64               // 主键
	ClusterId     int64               // 集群信息
	BuildNumber   int                 // 构建号
	Status        eumBuildStatus.Enum // 状态
	IsSuccess     bool                // 是否成功
	CreateAt      string              // 开始时间
	FinishAt      string              // 完成时间
	AppName       string              // 应用名称
	WorkflowsName string              // 工作流名称（文件的名称）
	BranchName    string              // 分支名称
}
