package response

// BuildManifestResponse 构建清单响应
type BuildManifestResponse struct {
	AppName       string // 应用名称
	GitName       string // 应用或库名称
	BuildNumber   int    // 构建号
	WorkflowsName string // 工作流名称
	DockerImage   string // 镜像名称
	GitBranch     string // GIT分支
	GitCommitId   string // git commitId
	CreateAt      string // 构建时间
}

// BuildManifestDetailResponse 构建清单详情响应
type BuildManifestDetailResponse struct {
	App          BuildManifestResponse   // 应用信息
	Dependencies []BuildManifestResponse // 依赖库信息
}
