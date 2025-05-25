package apps

import (
	"net/url"
	"path/filepath"
	"strings"

	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/utils/str"
)

// GitEO git仓库
type GitEO struct {
	Id       int    // 主键
	Name     string // Git名称
	Hub      string // git地址
	Branch   string // Git分支
	UserName string // 账户名称
	UserPwd  string // 账户密码
	Path     string // 存储目录
	IsApp    bool   // 是否为应用
}

func (receiver *GitEO) IsNil() bool {
	return receiver.Id == 0
}

// GetAbsolutePath 获取git存储的绝对路径 如："/var/lib/fops/git/fops/"
func (receiver *GitEO) GetAbsolutePath() string {
	return GitRoot + receiver.GetRelativePath()
}

// GetRelativePath 获取git存储的相对路径 如："fops/"
func (receiver *GitEO) GetRelativePath() string {
	if receiver.Path == "" || receiver.Path == "/" {
		receiver.Path = receiver.GetName()
	}
	// 移除前后/
	receiver.Path = strings.TrimPrefix(receiver.Path, "/")
	receiver.Path = strings.TrimSuffix(receiver.Path, "/")
	return receiver.Path + "/"
}

// GetName 获取仓库名称 如："fops"
func (receiver *GitEO) GetName() string {
	git := filepath.Base(receiver.Hub)
	return str.CutRight(git, ".git")
}

// GetAuthHub 获取带账号密码的地址 如："https://steden:123456@github.com/farseer-go/fs.git"
func (receiver *GitEO) GetAuthHub() string {
	parsedURL, err := url.Parse(receiver.Hub)
	exception.ThrowRefuseExceptionfBool(err != nil, "解析 URL 失败:%s", err)

	// 设置用户名和密码
	parsedURL.User = url.UserPassword(receiver.UserName, receiver.UserPwd)

	return parsedURL.String()
}

// GetRawContent 获取github仓库中的内容
func (receiver *GitEO) GetRawContent(filePath string) string {
	// 如："https://steden:123456@github.com/farseers/FOPS.git"
	gitUrl := receiver.GetAuthHub()
	if strings.Contains(gitUrl, "github.com") {
		// 移除.git后缀 https://raw.githubusercontent.com/farseers/FOPS/main/.fops/workflows/build.yml
		gitUrl = strings.TrimSuffix(gitUrl, ".git")
		gitUrl = gitUrl + "/" + receiver.Branch + "/" + filePath
		gitUrl = strings.ReplaceAll(gitUrl, "github.com", "raw.githubusercontent.com")
	}
	return gitUrl
}

// // GetWorkflowsRoot 获取工作流目录 如："/var/lib/fops/workflows/fops/"
// func (receiver *GitEO) GetWorkflowsRoot() string {
// 	return WorkflowsRoot + receiver.GetName() + "/"
// }
