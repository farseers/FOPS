package device

import (
	"context"
	"fops/domain/apps"
	"fops/domain/apps/event"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/file"
	"github.com/farseer-go/utils/str"
	"github.com/timandy/routine"
	"path/filepath"
	"sync"
)

func RegisterGitDevice() {
	container.Register(func() apps.IGitDevice { return &gitDevice{} })
}

type gitDevice struct {
}

func (device *gitDevice) GetGitPath(gitHub string) string {
	if gitHub == "" {
		return ""
	}
	var gitName = device.GetName(gitHub)
	return apps.GitRoot + gitName + "/"
}

func (*gitDevice) GetName(gitHub string) string {
	if gitHub == "" {
		return ""
	}
	git := filepath.Base(gitHub)
	return str.CutRight(git, ".git")
}

func (*gitDevice) RememberPassword(env apps.EnvVO, progress chan string) {
	exec.RunShell("git config --global credential.helper store", progress, env.ToMap(), "")
}

func (*gitDevice) ExistsGitProject(gitPath string) bool {
	// 如果Git存放的目录不存在，则创建
	if !file.IsExists(apps.GitRoot) {
		file.CreateDir766(apps.GitRoot)
	}
	return file.IsExists(gitPath)
}

func (device *gitDevice) Clear(git apps.GitEO, progress chan string) bool {
	// 获取Git存放的路径
	gitPath := git.GetAbsolutePath()
	exitCode := exec.RunShell("rm -rf "+gitPath, progress, nil, "")
	if exitCode != 0 {
		progress <- "Git清除失败"
		return false
	}
	return true
}

func (device *gitDevice) CloneOrPull(git apps.GitEO, progress chan string, ctx context.Context) bool {
	if progress == nil {
		progress = make(chan string, 100)
	}
	//progress <- "---------------------------------------------------------"

	// 先得到项目Git存放的物理路径
	gitPath := git.GetAbsolutePath()
	var execSuccess bool

	// 存在则使用pull
	if device.ExistsGitProject(gitPath) {
		//progress <- "开始拉取git " + git.Name + " 分支：" + git.Branch + " 仓库：" + git.Hub + "。"
		execSuccess = pull(gitPath, progress, ctx)
	} else {
		//progress <- "开始克隆git " + git.Name + " 分支：" + git.Branch + " 仓库：" + git.Hub + "。"
		execSuccess = device.clone(gitPath, git.GetAuthHub(), git.Branch, progress, ctx)
	}

	if execSuccess {
		// 通知更新git拉取时间
		event.GitCloneOrPulledEvent{GitId: git.Id}.PublishEvent()
	} else {
		progress <- "拉取出错了。"
	}
	return execSuccess
}

func (device *gitDevice) CloneOrPullAndDependent(lstGit []apps.GitEO, progress chan string, ctx context.Context) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始拉取git代码"
	var sw sync.WaitGroup
	result := true
	for _, git := range lstGit {
		sw.Add(1)
		g := git
		routine.Go(func() {
			defer sw.Done()
			if !device.CloneOrPull(g, progress, ctx) {
				result = false
			}
		})
	}
	sw.Wait()
	if result {
		progress <- "拉取完成。"
	}
	return result
}

func pull(savePath string, progress chan string, ctx context.Context) bool {
	exitCode := exec.RunShellContext(ctx, "git -C "+savePath+" pull --rebase", progress, nil, "")
	if exitCode != 0 {
		progress <- "Git拉取失败"
		return false
	}
	return true
}

func (device *gitDevice) clone(gitPath string, github string, branch string, progress chan string, ctx context.Context) bool {
	exitCode := exec.RunShellContext(ctx, "git clone -b "+branch+" "+github+" "+gitPath, progress, nil, "")
	if exitCode != 0 {
		progress <- "Git克隆失败"
		return false
	}
	return true
}
