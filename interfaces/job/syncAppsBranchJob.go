package job

import (
	"fops/domain/apps"
	"fops/domain/appsBranch"

	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/tasks"
)

// SyncAppsBranchJob 同步Git分支
func SyncAppsBranchJob(*tasks.TaskContext) {
	appsRepository := container.Resolve[apps.Repository]()
	appsBranchRepository := container.Resolve[appsBranch.Repository]()
	gitDevice := container.Resolve[apps.IGitDevice]()
	lstApp := appsRepository.ToUTList()
	lstUT := appsBranchRepository.ToList()

	// 读取远程分支，并更新状态
	lstApp.Foreach(func(appDO *apps.DomainObject) {
		// 在工作流根目录，获取远程分支名称
		lstBranch := gitDevice.GetRemoteBranch(fs.Context, appDO.GetWorkflowsRoot())
		if lstBranch.Count() == 0 {
			return
		}

		// 当前应用的本地分支
		lstLocalUT := lstUT.Where(func(item appsBranch.DomainObject) bool {
			return item.AppName == appDO.AppName
		}).ToList()

		lstBranch.Foreach(func(remoteBranchVO *apps.RemoteBranchVO) {
			// 找到数据库中的UT记录
			dbUT := lstLocalUT.Find(func(item *appsBranch.DomainObject) bool {
				return item.BranchName == remoteBranchVO.BranchName
			})
			// 不存在，则直接添加
			if dbUT == nil {
				do := appsBranch.DomainObject{AppName: appDO.AppName, BranchName: remoteBranchVO.BranchName, CommitId: remoteBranchVO.CommitId, CommitMessage: remoteBranchVO.CommitMessage, CommitAt: dateTime.Now(), BuildAt: dateTime.Now()}
				appsBranchRepository.Add(do)
				//lstUT.Add(do)
				return
			}

			// 不相等时，说明有新的提交，则替换
			if dbUT.CommitId != remoteBranchVO.CommitId {
				dbUT.CommitId = remoteBranchVO.CommitId
				dbUT.CommitMessage = remoteBranchVO.CommitMessage
				dbUT.CommitAt = dateTime.Now()
				dbUT.BuildErrorCount = 0
				dbUT.BuildAt = dateTime.Now()
				dbUT.BuildId = 0
				dbUT.BuildSuccess = false
				appsBranchRepository.UpdateByBranch(*dbUT)
			}
		})

		// 通过遍历本地分支，判断远程分支是否存在
		lstLocalUT.Foreach(func(utDO *appsBranch.DomainObject) {
			// 远程分支不存在，说明已经被删了
			if !lstBranch.Where(func(item apps.RemoteBranchVO) bool {
				return item.BranchName == utDO.BranchName
			}).Any() {
				flog.Infof("应用：%s的远程分支：%s不存在了，删除本地分支", utDO.AppName, utDO.BranchName)
				appsBranchRepository.DeleteBranch(utDO.AppName, utDO.BranchName)
			}
		})
	})
}
