package apps

import (
	"fmt"
	"os"
	"strconv"

	"github.com/farseer-go/fs/parse"
)

const (
	FopsRoot         = "/var/lib/fops/"                   // Fops根目录
	WithJsonPath     = "/var/lib/fops/dist/with.json"     // with.json文件位置
	KubeRoot         = "/var/lib/fops/kube/"              // kubectlConfig配置
	NpmModulesRoot   = "/var/lib/fops/npm"                // NpmModules
	DistRoot         = "/var/lib/fops/dist/"              // 编译保存的根目录
	GitRoot          = "/var/lib/fops/git/"               // GIT根目录
	DockerfilePath   = "/var/lib/fops/dist/Dockerfile"    // Dockerfile文件地址
	DockerIgnorePath = "/var/lib/fops/dist/.dockerignore" // DockerIgnore文件地址
	ShellRoot        = "/var/lib/fops/shell/"             // 生成Shell脚本的存放路径
	ActionsRoot      = "/var/lib/fops/actions/"           // 执行Actions的缓存目录
	WorkflowsRoot    = "/var/lib/fops/workflows/"         // 存放工作流文件的目录
	BackupRoot       = "/var/lib/fops/backup/"            // 存放备份文件的目录
)

// InitFopsDir 初始化目录
func InitFopsDir() {
	// sudo chmod -R 777 /var/lib/fops/
	_ = os.MkdirAll(FopsRoot, 0777)
	_ = os.MkdirAll(KubeRoot, 0777)
	_ = os.MkdirAll(NpmModulesRoot, 0777)
	_ = os.MkdirAll(DistRoot, 0777)
	_ = os.MkdirAll(GitRoot, 0777)
	_ = os.MkdirAll(ShellRoot, 0777)
	_ = os.MkdirAll(ActionsRoot, 0777)
	_ = os.MkdirAll(WorkflowsRoot, 0777)
	_ = os.MkdirAll(BackupRoot, 0777)
}

// EnvVO 构建时的环境变量
type EnvVO struct {
	BuildId     int64  // 构建主键
	ClusterId   int64  // 构建的集群ID
	BuildNumber int    // 构建版本号
	AppName     string // 项目名称
	AppGitRoot  string // Git仓库源代码根目录 /var/lib/fops/git/{gitName}/
	GitHub      string // Git仓库地址
	GitName     string // Git名称（项目的目录名称）
	DockerHub   string // Docker仓库地址
	DockerImage string // Docker镜像
	BranchName  string // 分支名称
	CommitId    string // 应用的CommitId
	Sha256sum   string // 应用的Sha256sum
}

// Print 打印环境变量
func (env *EnvVO) Print(progress chan string) {
	// 打印环境变量
	progress <- "---------------------------------------------------------"
	progress <- "环境变量："

	for k, v := range env.ToMap() {
		progress <- fmt.Sprint(k, "=", v)
	}
}

// ToMap 转成字典
func (env *EnvVO) ToMap() map[string]string {
	return map[string]string{
		"FopsRoot":       FopsRoot,
		"NpmModulesRoot": NpmModulesRoot,
		"DistRoot":       DistRoot,
		"KubeRoot":       KubeRoot,
		"Git_Root":       GitRoot,
		"Git_Hub":        env.GitHub,
		"Build_Id":       parse.ToString(env.BuildId),
		"Build_Number":   strconv.Itoa(env.BuildNumber),
		"App_Name":       env.AppName,
		"App_GitRoot":    env.AppGitRoot,
		//"App_Domain":     env.ProjectDomain,
		//"App_EntryPoint": env.ProjectEntryPoint,
		//"App_EntryPort":  strconv.Itoa(env.ProjectEntryPort),
		"Docker_Hub":   env.DockerHub,
		"Docker_Image": env.DockerImage,
		"Git_Name":     env.GitName,
	}
}
