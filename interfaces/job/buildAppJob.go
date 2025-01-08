package job

import (
	"fops/domain/_/eumBuildType"
	"fops/domain/apps"

	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/tasks"
)

// BuildAppJob 开启构建
func BuildAppJob(*tasks.TaskContext) {
	appsRepository := container.Resolve[apps.Repository]()
	buildEO := appsRepository.GetUnBuildInfo(eumBuildType.Manual)
	if buildEO.IsNil() {
		if traceContext := trace.CurTraceContext.Get(); traceContext != nil {
			traceContext.Ignore()
		}
		return
	}

	// 应用不存在
	if appDO := appsRepository.ToEntity(buildEO.AppName); appDO.IsNil() {
		appsRepository.SetCancel(buildEO.Id, apps.EnvVO{})
		return
	}

	// 设置为构建中
	appsRepository.SetBuilding(buildEO.Id)

	// 开始构建
	buildEO.StartBuild()
}
