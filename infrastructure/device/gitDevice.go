package device

import (
	"context"
	"fmt"
	"fops/domain/apps"

	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/file"
)

func RegisterGitDevice() {
	container.Register(func() apps.IGitDevice { return &gitDevice{} })
}

type gitDevice struct {
}

func (receiver *gitDevice) PullWorkflows(ctx context.Context, gitPath, branch string, gitRemote string, progress chan string) bool {
	if !file.IsExists(gitPath) {
		file.CreateDir766(gitPath)
		exec.RunShell("git init", progress, nil, gitPath, true)
		exec.RunShell(fmt.Sprintf("git remote add -f origin %s", gitRemote), progress, nil, gitPath, true)
		exec.RunShell("git config core.sparsecheckout true", progress, nil, gitPath, true)
		exec.RunShell("echo .fops/workflows/ >> .git/info/sparse-checkout", progress, nil, gitPath, true)
	}

	var exitCode int
	for i := 0; i < 3; i++ {
		select {
		case <-ctx.Done():
			progress <- "同步工作流文件失败，停止构建"
			return false
		default:
			if exitCode = exec.RunShellContext(ctx, fmt.Sprintf("git pull origin %s", branch), progress, nil, gitPath, true); exitCode == 0 {
				return true
			}
			if i == 2 {
				progress <- "同步工作流文件失败，停止构建"
			} else {
				progress <- "拉取失败，正在尝试重新拉取"
			}
		}
	}
	return exitCode == 0
}
