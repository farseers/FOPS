package job

import (
	"fops/domain/_/eumBuildStatus"
	"fops/domain/_/eumBuildType"
	"fops/domain/apps"
	"fops/domain/appsBranch"

	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/tasks"
)

// AutoBuildAppJob 自动构建
func AutoBuildAppJob(*tasks.TaskContext) {
	appsRepository := container.Resolve[apps.Repository]()
	appsBranchRepository := container.Resolve[appsBranch.Repository]()
	appsBranchDO := appsBranchRepository.GetUnRunUT()
	if appsBranchDO.IsNil() {
		return
	}

	appDO := appsRepository.ToEntity(appsBranchDO.AppName)
	// 应用不存在
	if appDO.IsNil() {
		flog.Warningf("自动构建时，应用名称不存在：%s", appsBranchDO.AppName)
		return
	}

	buildDO := apps.BuildEO{
		BuildType:     eumBuildType.Auto,
		BuildServerId: core.AppId,
		BuildNumber:   0,
		CreateAt:      dateTime.Now(),
		FinishAt:      dateTime.Now(),
		Env:           apps.EnvVO{},
		AppName:       appsBranchDO.AppName,
		WorkflowsName: appDO.UTWorkflowsName,
		BranchName:    appsBranchDO.BranchName,
		Status:        eumBuildStatus.Building,
	}
	err := appsRepository.AddBuild(&buildDO)
	exception.ThrowRefuseExceptionError(err)

	// 启动构建
	buildDO.StartBuild()

	// 更新状态
	appsBranchDO.BuildId = buildDO.Id
	appsBranchDO.BuildSuccess = buildDO.IsSuccess
	appsBranchDO.BuildAt = dateTime.Now()
	if !buildDO.IsSuccess {
		appsBranchDO.BuildErrorCount++
	}
	// 如果自动构建被取消了，则全部暂停
	if buildDO.Status == eumBuildStatus.Cancel {
		appsBranchDO.BuildErrorCount = 3
	}
	appsBranchRepository.UpdateByBranch(appsBranchDO)
}
