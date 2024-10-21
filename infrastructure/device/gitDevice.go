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
		exec.RunShellContext(ctx, "git init", progress, nil, gitPath, true)
		exec.RunShellContext(ctx, fmt.Sprintf("git remote add -f origin %s", gitRemote), progress, nil, gitPath, true)
		exec.RunShellContext(ctx, "git config core.sparsecheckout true", progress, nil, gitPath, true)
		exec.RunShellContext(ctx, "echo .fops/workflows/ >> .git/info/sparse-checkout", progress, nil, gitPath, true)
	}

	var exitCode int
	for i := 0; i < 3; i++ {
		if exitCode = exec.RunShell(fmt.Sprintf("git pull origin %s", branch), progress, nil, gitPath, true); exitCode == 0 {
			break
		}
		if i == 2 {
			progress <- "同步工作流文件失败，停止构建"
		} else {
			progress <- "拉取失败，正在尝试重新拉取"
		}

	}
	return exitCode == 0
}
