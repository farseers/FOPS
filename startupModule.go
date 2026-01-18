package main

import (
	"fops/domain/apps"
	configure2 "fops/domain/configure"
	"fops/infrastructure"
	"fops/interfaces"
	"os"
	"runtime/debug"

	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/modules"
	"github.com/farseer-go/utils/exec"
)

type StartupModule struct {
}

func (module StartupModule) DependsModule() []modules.FarseerModule {
	return []modules.FarseerModule{infrastructure.Module{}, interfaces.Module{}}
}

func (module StartupModule) PreInitialize() {
	debug.SetGCPercent(50)
}

func (module StartupModule) Initialize() {
}

func (module StartupModule) PostInitialize() {
	// 替换Fops.Proxy的值
	if configure.GetString("Fops.Proxy") == "global.proxy" {
		configureDO := container.Resolve[configure2.Repository]().ToEntityByKey("global", "proxy")
		_ = os.Setenv("Fops_Proxy", configureDO.Value)
		flog.Infof("使用配置管理global.proxy的代理设置：%s", configureDO.Value)
	}

	// 使用git代理
	if proxyAgent := configure.GetString("Fops.Proxy"); proxyAgent != "" {
		flog.Info("开启Git代理：", proxyAgent)
		exec.RunShellCommand("git config --global http.https://github.com.proxy "+proxyAgent, nil, "", true)
		exec.RunShellCommand("git config --global https.https://github.com.proxy "+proxyAgent, nil, "", true)
	} else {
		flog.Info("未开启Git代理")
		exec.RunShellCommand("git config --global --unset http.https://github.com.proxy", nil, "", false)
		exec.RunShellCommand("git config --global --unset https.https://github.com.proxy", nil, "", false)
	}

	// 初始化目录
	apps.InitFopsDir()
}

func (module StartupModule) Shutdown() {
}
