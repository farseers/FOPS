package device

import (
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

func (receiver *gitDevice) PullWorkflows(gitPath, branch string, gitRemote string, progress chan string) bool {
	if !file.IsExists(gitPath) {
		file.CreateDir766(gitPath)
		exec.RunShell("git init", progress, nil, gitPath, true)
		exec.RunShell(fmt.Sprintf("git remote add -f %s %s", branch, gitRemote), progress, nil, gitPath, true)
		exec.RunShell("git config core.sparsecheckout true", progress, nil, gitPath, true)
		exec.RunShell("echo .fops/workflows/ >> .git/info/sparse-checkout", progress, nil, gitPath, true)
	}

	var exitCode int
	for i := 0; i < 3; i++ {
		if exitCode = exec.RunShell(fmt.Sprintf("git pull %s %s", branch, branch), progress, nil, gitPath, true); exitCode == 0 {
			break
		}
		progress <- "同步工作流文件失败"
	}
	return exitCode == 0
}