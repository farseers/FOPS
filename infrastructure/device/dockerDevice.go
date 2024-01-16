package device

import (
	"context"
	"fmt"
	"fops/domain/apps"
	"fops/domain/apps/event"
	"fops/domain/cluster"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/file"
	"github.com/farseer-go/utils/str"
	"os"
	"regexp"
	"strconv"
)

func RegisterDockerDevice() {
	container.Register(func() apps.IDockerDevice { return &dockerDevice{} })
}

type dockerDevice struct {
}

func (dockerDevice) GetDockerHub(dockerHubAddress string) string {
	var dockerHub = "localhost"
	if dockerHubAddress != "" {
		dockerHub = dockerHubAddress
		dockerHub = str.CutRight(dockerHub, "/")
	}
	return dockerHub
}

func (device dockerDevice) GetDockerImage(dockerHubAddress string, projectName string, buildNumber int) string {
	return device.GetDockerHub(dockerHubAddress) + "/" + projectName + ":" + strconv.Itoa(buildNumber)
}

func (dockerDevice) Login(dockerHub string, loginName string, loginPwd string, progress chan string, env apps.EnvVO, ctx context.Context) bool {
	progress <- "---------------------------------------------------------"
	if dockerHub != "" && loginName != "" {
		var result = exec.RunShellContext(ctx, "docker login "+dockerHub+" -u "+loginName+" -p "+loginPwd, progress, env.ToMap(), "")
		if result != 0 {
			progress <- "镜像仓库登陆失败。"
			return false
		}
	}

	progress <- "镜像仓库登陆成功。"
	return true
}

func (dockerDevice) ExistsDockerfile(dockerfilePath string) bool {
	return file.IsExists(dockerfilePath)
}

func (dockerDevice) CreateDockerfile(projectName string, dockerfileContent string, ctx context.Context) {
	if file.IsExists(apps.DockerfilePath) {
		_ = os.RemoveAll(apps.DockerfilePath)
	}
	file.WriteString(apps.DockerfilePath, dockerfileContent)
}

func (dockerDevice) Build(env apps.EnvVO, progress chan string, ctx context.Context) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始镜像打包。"

	// 打包
	var result = exec.RunShellContext(ctx, "docker build -t "+env.DockerImage+" --network=host -f "+apps.DockerfilePath+" "+apps.DistRoot, progress, env.ToMap(), apps.DistRoot)
	if result == 0 {
		progress <- "镜像打包完成。"
	} else {
		progress <- "镜像打包出错了。"
	}
	return result == 0
}

func (dockerDevice) Push(env apps.EnvVO, progress chan string, ctx context.Context) bool {
	defer func() {
		// 上传完后，删除本地镜像
		exec.RunShellContext(ctx, "docker rmi "+env.DockerImage, progress, env.ToMap(), "")
	}()

	// 上传
	var result = exec.RunShellContext(ctx, "docker push "+env.DockerImage, progress, env.ToMap(), "")

	if result == 0 {
		progress <- "镜像上传完成。"

		// 上传成功后，需要更新项目中的镜像版本属性
		event.DockerPushedEvent{
			BuildNumber: env.BuildNumber,
			AppName:     env.AppName,
			ImageName:   env.DockerImage,
		}.PublishEvent()
		return true
	}

	progress <- "镜像上传出错了。"
	return false
}

func (dockerDevice) SetImages(cluster cluster.DomainObject, appName string, dockerImages string, progress chan string, ctx context.Context) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始更新Docker Swarm的镜像版本。"

	var exitCode = exec.RunShellContext(ctx, fmt.Sprintf("docker service update --image %s %s", dockerImages, appName), progress, nil, "")
	if exitCode != 0 {
		progress <- "Docker Swarm更新镜像失败。"
		return false
	}
	progress <- "Docker Swarm更新镜像版本完成。"
	return true
}

func (dockerDevice) ExistsDocker(cluster cluster.DomainObject, appName string) bool {
	progress := make(chan string, 1000)
	// docker service inspect fops
	var exitCode = exec.RunShell(fmt.Sprintf("docker service inspect %s", appName), progress, nil, "")
	lst := collections.NewListFromChan(progress)
	if exitCode != 0 {
		if lst.Contains("[]") && lst.ContainsPrefix("Status: Error: no such service:") {
			return false
		}
		exception.ThrowWebException(403, "获取应用信息时失败。")
		return false
	}
	if lst.Contains("[]") && lst.ContainsPrefix("Status: Error: no such service:") {
		return false
	}
	return lst.ContainsAny(fmt.Sprintf("\"Name\": \"%s\"", appName))
}

func (dockerDevice) CreateService(appName, dockerNodeRole, additionalScripts, dockerNetwork string, dockerReplicas int, dockerImages string, progress chan string, ctx context.Context) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始创建Docker Swarm容器服务。"

	shell := fmt.Sprintf("docker service create --name %s --replicas %v -d --network=%s --constraint node.role==%s --mount type=bind,src=/etc/localtime,dst=/etc/localtime %s %s", appName, dockerReplicas, dockerNetwork, dockerNodeRole, additionalScripts, dockerImages)
	var exitCode = exec.RunShellContext(ctx, shell, progress, nil, "")
	if exitCode != 0 {
		progress <- "创建Docker Swarm容器失败了。"
		return false
	}
	return true
}

func (dockerDevice) DeleteService(appName string, progress chan string) bool {
	// docker service rm fops
	var exitCode = exec.RunShell(fmt.Sprintf("docker service rm %s", appName), progress, nil, "")
	if exitCode != 0 {
		progress <- "删除Docker Swarm容器失败了。"
		return false
	}
	return true
}

func (dockerDevice) SetReplicas(cluster cluster.DomainObject, appName string, dockerReplicas int, progress chan string) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始更新Docker Swarm的副本数量。"

	var exitCode = exec.RunShell(fmt.Sprintf("docker service update --replicas %v %s", dockerReplicas, appName), progress, nil, "")
	if exitCode != 0 {
		progress <- "Docker Swarm的副本数量更新失败。"
		return false
	}
	progress <- "Docker Swarm的副本数量更新完成。"
	return true
}

// ClearImages 清除镜像
func (dockerDevice) ClearImages(progress chan string) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始清除镜像。"

	var exitCode = exec.RunShell("docker system prune -f && docker builder prune && docker rmi $(docker images -f \"dangling=true\" -q)", progress, nil, "")
	if exitCode != 0 {
		progress <- "清除镜像镜像失败。"
		return false
	}
	progress <- "清除镜像完成。"
	return true
}

func (dockerDevice) Restart(cluster cluster.DomainObject, appName string, progress chan string) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始重启Docker Swarm的容器。"

	var exitCode = exec.RunShell(fmt.Sprintf("docker service update --force %s", appName), progress, nil, "")
	if exitCode != 0 {
		progress <- "Docker Swarm的容器重启失败。"
		return false
	}
	progress <- "Docker Swarm的容器重启完成。"
	return true
}

func (dockerDevice) GetVersion() string {
	receiveOutput := make(chan string, 100)
	exec.RunShell("docker version --format '{{.Server.Version}}'", receiveOutput, nil, "")
	lst := collections.NewListFromChan(receiveOutput)
	re := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	for _, s := range lst.ToArray() {
		if re.MatchString(s) {
			return s
		}
	}
	return ""
}
