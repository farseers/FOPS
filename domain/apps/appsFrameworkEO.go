package apps

import "regexp"

var commitIdRegexp = regexp.MustCompile(`^[0-9a-fA-F]{7,40}$`)

// AppsFrameworkEO 应用与框架关系实体
type AppsFrameworkEO struct {
	AppName     string // 应用名称
	FrameworkId int64  // 框架ID
	CommitId    string // 框架提交ID
}

// 判断是分支名称还是CommitId
func (receiver *AppsFrameworkEO) IsCommitId() bool {
	return commitIdRegexp.MatchString(receiver.CommitId)
}
