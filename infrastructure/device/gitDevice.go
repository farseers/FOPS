package device

import (
	"context"
	"encoding/json"
	"fmt"
	"fops/domain/apps"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/file"
	"github.com/farseer-go/utils/http"
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

func (receiver *gitDevice) GetLocalBranch(ctx context.Context, gitPath string) collections.List[apps.RemoteBranchVO] {
	lst := collections.NewList[apps.RemoteBranchVO]()
	wait := exec.RunShellContext(ctx, "git", []string{"ls-remote", "--heads", "origin"}, nil, gitPath, false)
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

func (receiver *gitDevice) GetRemoteBranch(ctx context.Context, gitAuthHb string) collections.List[apps.RemoteBranchVO] {
	lst := collections.NewList[apps.RemoteBranchVO]()
	wait := exec.RunShellContext(ctx, "git", []string{"ls-remote", gitAuthHb}, nil, "", false)
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

// CreateTag 核心方法：使用自定义 HTTP 工具实现打 Tag
func (receiver *gitDevice) CreateTag(ctx context.Context, gitAuthHb, branchOrCommitId, tagName string) error {
	// 1. 解析 URL 和 Token
	u, err := url.Parse(gitAuthHb)
	if err != nil {
		return fmt.Errorf("gitAuthHb 解析失败: %w", err)
	}
	if u.User == nil {
		return fmt.Errorf("gitAuthHb 中缺少认证信息")
	}

	// 兼容 user:token 和 token@host 两种格式
	token, ok := u.User.Password()
	if !ok {
		token = u.User.Username()
	}

	// 2. 构造 API 基础地址 (支持 GitHub 和 GHE)
	repoPath := strings.TrimSuffix(strings.TrimPrefix(u.Path, "/"), ".git")
	apiBase := "https://api.github.com"
	if u.Host != "github.com" {
		apiBase = fmt.Sprintf("https://%s/api/v3", u.Host)
	}

	// 3. 准备通用 Headers
	headers := map[string]any{
		"Authorization": "Bearer " + token,
		"Accept":        "application/vnd.github+json",
	}

	// 4. 获取最新 SHA
	shaURL := fmt.Sprintf("%s/repos/%s/commits/%s", apiBase, repoPath, branchOrCommitId)
	respBody, statusCode, _, err := http.RequestProxyConfigure("GET", shaURL, headers, nil, "", 30)
	if err != nil {
		return fmt.Errorf("获取 SHA 请求失败: %w", err)
	}
	if statusCode != 200 {
		return fmt.Errorf("获取 SHA 失败 (状态码 %d): %s", statusCode, respBody)
	}

	// 解析 SHA
	var result struct {
		SHA string `json:"sha"`
	}
	if err := json.Unmarshal([]byte(respBody), &result); err != nil {
		return fmt.Errorf("解析 SHA 响应失败: %w", err)
	}

	// 5. 创建 Tag
	payload := map[string]string{
		"ref": "refs/tags/" + tagName,
		"sha": result.SHA,
	}

	refURL := fmt.Sprintf("%s/repos/%s/git/refs", apiBase, repoPath)
	respBody, statusCode, _, err = http.RequestProxyConfigure("POST", refURL, headers, payload, "application/json", 30)
	if err != nil {
		return fmt.Errorf("创建标签请求失败: %w", err)
	}
	if statusCode != 201 {
		return fmt.Errorf("创建标签失败 (状态码 %d): %s", statusCode, respBody)
	}

	return nil
}
