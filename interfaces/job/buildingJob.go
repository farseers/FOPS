package job

import (
	"fops/domain/apps"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/tasks"
)

// BuildingJob 开启构建
func BuildingJob(*tasks.TaskContext) {
	appsRepository := container.Resolve[apps.Repository]()
	buildEO := appsRepository.GetUnBuildInfo()
	if traceContext := trace.CurTraceContext.Get(); traceContext != nil && buildEO.IsNil() {
		traceContext.Ignore()
		return
	}

	appDO := appsRepository.ToEntity(buildEO.AppName)
	// 应用不存在
	if appDO.IsNil() {
		appsRepository.SetCancel(buildEO.Id, apps.EnvVO{}, nil)
		return
	}

	// 设置为构建中
	appsRepository.SetBuilding(buildEO.Id)

	// 开始构建
	buildEO.StartBuild()
}
