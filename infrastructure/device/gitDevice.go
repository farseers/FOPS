package device

import (
	"context"
	"fmt"
	"fops/domain/apps"
	"strings"

	"github.com/farseer-go/collections"
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
	exec.RunShell("git config --global http.timeout 10", progress, nil, "", true)

	var exitCode int
	for i := 0; i < 3; i++ {
		select {
		case <-ctx.Done():
			progress <- "同步工作流文件失败，停止构建"
			return false
		default:
			if exitCode = exec.RunShellContext(ctx, fmt.Sprintf("timeout 10 git pull origin %s", branch), progress, nil, gitPath, true); exitCode == 0 {
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

func (receiver *gitDevice) GetRemoteBranch(ctx context.Context, gitPath string) collections.List[apps.RemoteBranchVO] {
	lst := collections.NewList[apps.RemoteBranchVO]()
	progress := make(chan string, 10000)
	// git ls-remote --heads
	// git branch -vr
	if exitCode := exec.RunShellContext(ctx, "git remote update origin --prune && timeout 10 git ls-remote --heads", progress, nil, gitPath, false); exitCode != 0 {
		return lst
	}
	lstContent := collections.NewListFromChan(progress)
	if lstContent.Any() {
		lstContent.RemoveAt(0)
	}
	for _, content := range lstContent.ToArray() {
		fields := strings.Fields(content)
		if len(fields) < 2 {
			continue
		}
		if len(fields[0]) < 16 {
			continue
		}
		remoteBranch := apps.RemoteBranchVO{
			CommitId: fields[0][:16],
			//CommitMessage: content[len(fields[0])+len(fields[1])+3:], // 消息带有空格，不能直接取fields[2]
		}
		remoteBranch.BranchName, _ = strings.CutPrefix(fields[1], "refs/heads/")
		lst.Add(remoteBranch)
	}
	return lst
}
