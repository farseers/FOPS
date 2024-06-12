package apps

import (
	"fops/domain/_/eumBuildStatus"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"time"
)

// Repository 仓储接口
type Repository interface {
	data.IRepository[DomainObject] // IRepository 通用的仓储接口
	UpdateApp(do DomainObject) error
	UpdateDockerVer(appName string, dockerVer int, imageName string) (int64, error)        // UpdateDockerVer 修改镜像版本
	UpdateClusterVer(appName string, dicClusterVer map[int64]*ClusterVerVO) (int64, error) // UpdateClusterVer 修改集群的镜像版本
	buildRepository
	gitRepository
}

type buildRepository interface {
	GetBuildNumber(appName string) int                                                     // 获取构建的编号
	AddBuild(eo BuildEO) error                                                             // 添加构建
	ToBuildList(appName string, pageSize int, pageIndex int) collections.PageList[BuildEO] // 查询构建列表
	GetUnBuildInfo() BuildEO                                                               // 获取未构建的任务
	SetBuilding(id int64)                                                                  // 设置任务为构建中
	SetSuccess(id int64)                                                                   // Success 任务完成
	SetCancel(id int64)                                                                    // Cancel 主动取消任务
	GetStatus(id int64) eumBuildStatus.Enum                                                // GetStatus 获取构建状态
}

type gitRepository interface {
	ToGitEntity(id int64) GitEO
	ToGitList(lstIds collections.List[int64]) collections.List[GitEO]
	ToGitListAll(isApp int) collections.List[GitEO]
	AddGit(eo GitEO) error
	UpdateGit(eo GitEO) (int64, error)
	DeleteGit(id int64) (int64, error)
	ExistsGit(id int64) bool
	UpdateForTime(id int, pullAt time.Time) (int64, error) // 修改GIT的拉取时间
}
