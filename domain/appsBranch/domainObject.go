package appsBranch

import "github.com/farseer-go/fs/dateTime"

// DomainObject 应用分支
type DomainObject struct {
	AppName         string            // 应用名称
	BranchName      string            // 分支名称
	CommitId        string            // 当前分支最后提交ID
	CommitMessage   string            // 提交消息
	CommitAt        dateTime.DateTime // 同步时间
	BuildId         int64             // 对应的构建ID
	BuildSuccess    bool              // 是否构建成功
	BuildErrorCount int               // 构建失败次数
	BuildAt         dateTime.DateTime // 构建时间
}

func (receiver DomainObject) IsNil() bool {
	return receiver.AppName == "" && receiver.BranchName == ""
}
