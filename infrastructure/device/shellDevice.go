package device

import (
	"context"
	"fops/domain/apps"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/file"
)

func RegisterShellDevice() {
	container.Register(func() apps.IShellDevice { return &shellDevice{} })
}

type shellDevice struct {
}

func (shellDevice) ExecShell(env apps.EnvVO, shellScript string, progress chan string, ctx context.Context) bool {
	// 每次执行时，需要生成shell脚本
	path := apps.ShellRoot + "fops_" + parse.ToString(env.BuildId) + ".sh"
	file.WriteString(path, shellScript)

	// 执行脚本
	var exitCode = exec.RunShellContext(ctx, "/bin/sh -xe {path}", progress, env.ToMap(), apps.DistRoot, true)
	if exitCode == 0 {
		progress <- "执行脚本完成。"
	} else {

		progress <- "执行脚本出错了。"
	}
	return exitCode == 0
}
