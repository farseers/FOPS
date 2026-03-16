package device

import (
	"context"
	"fops/domain/apps"
	"os"
	"path/filepath"
	"strings"
	"time"

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

		// git init
		wait := exec.RunShellContext(ctx, "git", []string{"init"}, nil, gitPath, true)
		if wait.WaitToChan(progress) != 0 {
			progress <- "git init 失败"
			return false
		}

		// git remote add
		wait = exec.RunShellContext(ctx, "git", []string{"remote", "add", "-f", "origin", gitRemote}, nil, gitPath, true)
		if wait.WaitToChan(progress) != 0 {
			progress <- "添加远程仓库失败"
			return false
		}

		// 配置稀疏检出
		wait = exec.RunShellContext(ctx, "git", []string{"config", "core.sparsecheckout", "true"}, nil, gitPath, true)
		wait.WaitToChan(progress)

		// 写入稀疏检出配置（覆盖写入，避免重复）
		sparseCheckoutPath := filepath.Join(gitPath, ".git", "info", "sparse-checkout")
		sparseCheckoutDir := filepath.Dir(sparseCheckoutPath)
		if !file.IsExists(sparseCheckoutDir) {
			file.CreateDir766(sparseCheckoutDir)
		}
		// 使用 WriteFile 覆盖写入，避免重复追加
		if err := os.WriteFile(sparseCheckoutPath, []byte(".fops/workflows/\n"), 0644); err != nil {
			progress <- "写入稀疏检出配置失败: " + err.Error()
			return false
		}
	}

	// 使用本地配置，不污染全局
	wait := exec.RunShellContext(ctx, "git", []string{"config", "http.timeout", "10"}, nil, gitPath, true)
	wait.WaitToChan(progress)

	// 先切换到目标分支
	// 尝试检出分支，如果不存在则创建并跟踪远程分支
	wait = exec.RunShellContext(ctx, "git", []string{"checkout", "-B", branch, "origin/" + branch}, nil, gitPath, true)
	if wait.WaitToChan(progress) != 0 {
		progress <- "切换到分支 " + branch + " 失败"
		return false
	}

	var exitCode int
	for i := 0; i < 3; i++ {
		select {
		case <-ctx.Done():
			progress <- "同步工作流文件失败，停止构建"
			return false
		default:
			// 使用 context 控制超时，不依赖外部 timeout 命令
			pullCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			wait := exec.RunShellContext(pullCtx, "git", []string{"pull", "origin", branch}, nil, gitPath, true)
			if exitCode = wait.WaitToChan(progress); exitCode == 0 {
				return true
			}
		}
	}
	return exitCode == 0
}

func (receiver *gitDevice) GetRemoteBranch(ctx context.Context, gitPath string) collections.List[apps.RemoteBranchVO] {
	lst := collections.NewList[apps.RemoteBranchVO]()
	wait := exec.RunShellContext(ctx, "timeout", []string{"10", "git", "ls-remote", "--heads", "origin"}, nil, gitPath, false)
	lstContent, exitCode := wait.WaitToList()
	if exitCode != 0 {
		return lst
	}

	for _, content := range lstContent.ToArray() {
		// 跳过空行
		content = strings.TrimSpace(content)
		if content == "" {
			continue
		}

		fields := strings.Fields(content)
		if len(fields) < 2 {
			continue
		}

		commitHash := fields[0]
		if len(commitHash) < 16 {
			continue
		}

		// refs/heads/xxx -> xxx
		branchName, found := strings.CutPrefix(fields[1], "refs/heads/")
		if !found {
			continue
		}

		remoteBranch := apps.RemoteBranchVO{
			CommitId: commitHash[:16],
		}
		remoteBranch.BranchName = branchName
		lst.Add(remoteBranch)
	}
	return lst
}
