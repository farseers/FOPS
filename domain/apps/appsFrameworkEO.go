package apps

// AppsFrameworkEO 应用与框架关系实体
type AppsFrameworkEO struct {
	AppName      string // 应用名称
	FrameworkId  int64  // 框架ID
	CommitId     string // 框架提交ID
	IsAutoUpdate bool   // 是否自动更新(false:锁定;true:永远拉main)
}
