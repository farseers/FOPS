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
		lstResult, wait := exec.RunShellContext(ctx, "git", []string{"init"}, nil, gitPath, true)
		if exitCode := exec.SaveToChan(progress, lstResult, wait); exitCode != 0 {
			progress <- "git init 失败"
			return false
		}

		// git remote add
		lstResult, wait = exec.RunShellContext(ctx, "git", []string{"remote", "add", "-f", "origin", gitRemote}, nil, gitPath, true)
		if exitCode := exec.SaveToChan(progress, lstResult, wait); exitCode != 0 {
			progress <- "添加远程仓库失败"
			return false
		}

		// 配置稀疏检出
		lstResult, wait = exec.RunShellContext(ctx, "git", []string{"config", "core.sparsecheckout", "true"}, nil, gitPath, true)
		exec.SaveToChan(progress, lstResult, wait)

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
	lstResult, wait := exec.RunShellContext(ctx, "git", []string{"config", "http.timeout", "10"}, nil, gitPath, true)
	exec.SaveToChan(progress, lstResult, wait)

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

			lstResult, wait := exec.RunShellContext(pullCtx, "timeout", []string{"10", "git", "pull", "origin", branch}, nil, gitPath, true)
			if exitCode = exec.SaveToChan(progress, lstResult, wait); exitCode == 0 {
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
	lstResult, wait := exec.RunShellContext(ctx, "timeout", []string{"10", "git", "ls-remote", "--heads", "origin"}, nil, gitPath, false)
	if wait() != 0 {
		return lst
	}

	lstContent := collections.NewListFromChan(lstResult)
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
		if len(commitHash) < 8 {
			continue
		}

		// refs/heads/xxx -> xxx
		branchName, found := strings.CutPrefix(fields[1], "refs/heads/")
		if !found {
			continue
		}

		remoteBranch := apps.RemoteBranchVO{
			CommitId: commitHash[:8],
		}
		remoteBranch.BranchName = branchName
		lst.Add(remoteBranch)
	}
	return lst
}
