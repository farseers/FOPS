package apps

import "github.com/farseer-go/fs/dateTime"

// BuildManifestEO 构建清单实体
type BuildManifestEO struct {
	AppName       string            // 应用名称 (主键1)
	GitName       string            // 应用或库名称 (主键2)
	BuildNumber   int               // 构建号
	WorkflowsName string            // 工作流名称
	DockerImage   string            // 镜像名称
	GitId         int               // Git主键
	GitBranch     string            // GIT分支
	GitCommitId   string            // git commitId
	CreateAt      dateTime.DateTime // 构建时间
}

func (receiver *BuildManifestEO) IsNil() bool {
	return receiver.AppName == "" && receiver.GitName == ""
}
