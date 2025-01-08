package apps

import (
	"fops/domain/_/eumBuildStatus"
	"fops/domain/_/eumBuildType"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
)

// Repository 仓储接口
type Repository interface {
	data.IRepository[DomainObject] // IRepository 通用的仓储接口
	UpdateApp(do DomainObject) error
	UpdateDockerVer(appName string, dockerVer int, imageName string) (int64, error)                            // UpdateDockerVer 修改镜像版本
	UpdateClusterVer(appName string, dicClusterVer collections.Dictionary[int64, ClusterVerVO]) (int64, error) // UpdateClusterVer 修改集群的镜像版本
	UpdateInspect(lst collections.List[DomainObject]) (int64, error)                                           // UpdateInspect 更新从集群中获取到的实例、副本数量
	ToListBySys(isSys bool) collections.List[DomainObject]
	ToShortList(isAll bool) collections.List[ShortEO]
	ToUTList() collections.List[DomainObject] // 获取需要UT的应用
	buildRepository
	gitRepository
}

type buildRepository interface {
	GetBuildNumber(appName string) int                                                     // 获取构建的编号
	AddBuild(eo *BuildEO) error                                                            // 添加构建
	ToBuildList(appName string, pageSize int, pageIndex int) collections.PageList[BuildEO] // 查询构建列表
	GetUnBuildInfo(buildType eumBuildType.Enum) BuildEO                                    // 获取未构建的任务
	SetBuilding(id int64)                                                                  // 设置任务为构建中
	UpdateBuilding(id int64, env EnvVO)                                                    // 更新构建任务
	SetSuccess(id int64, env EnvVO)                                                        // Success 任务完成
	SetSuccessForFops(id int64)                                                            // 设置任务为构建成功
	SetCancel(id int64, env EnvVO)                                                         // Cancel 主动取消任务
	SetFail(id int64, env EnvVO)                                                           // SetFail 任务失败
	GetStatus(id int64) eumBuildStatus.Enum                                                // GetStatus 获取构建状态
	UpdateFailDockerImage(appName string, dockerImage string) (int64, error)               // UpdateFailDockerImage 更新构建中状态的构建记录
	GetLastBuilding() BuildEO                                                              // 获取最后一次构建
	ToBuildEntity(id int64) BuildEO                                                        // 获取构建对象
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
