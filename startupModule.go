package main

import (
	"fops/domain/apps"
	"fops/infrastructure"
	"fops/interfaces"
	"github.com/farseer-go/fs/configure"
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
}

func (module StartupModule) Initialize() {
}

func (module StartupModule) PostInitialize() {
	// 当前主机系统不需要设置了。改到工作流中执行
	// 使用git代理
	receiveOutput := make(chan string, 100)
	if gitAgent := configure.GetString("Fops.GitAgent"); gitAgent != "" {
		flog.Info("开启Git代理：", gitAgent)
		exec.RunShell("git config --global http.https://github.com.proxy "+gitAgent, receiveOutput, nil, "", true)
		exec.RunShell("git config --global https.https://github.com.proxy "+gitAgent, receiveOutput, nil, "", true)
	} else {
		exec.RunShell("git config --global --unset http.https://github.com.proxy", receiveOutput, nil, "", false)
		exec.RunShell("git config --global --unset https.https://github.com.proxy", receiveOutput, nil, "", false)
	}

	// 初始化目录
	apps.InitFopsDir()

	//appsApp.BuildAdd("lbn", 1, "本地", container.Resolve[apps.Repository](), container.Resolve[cluster.Repository]())
}

func (module StartupModule) Shutdown() {
}
