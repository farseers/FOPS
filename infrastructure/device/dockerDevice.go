package device

import (
	"bytes"
	"context"
	"fmt"
	"fops/domain/apps"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/utils/exec"
	"github.com/farseer-go/utils/str"
	"path"
	"regexp"
	"strconv"
	"strings"
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

func (device dockerDevice) GetDockerImage(dockerHubAddress string, appName string, buildNumber int) string {
	return device.GetDockerHub(dockerHubAddress) + ":" + appName + "." + strconv.Itoa(buildNumber)
}

func (dockerDevice) Run(dockerName string, network string, dockerImage string, args []string, useRm bool, env apps.EnvVO, progress chan string, ctx context.Context) bool {
	bf := bytes.Buffer{}
	bf.WriteString("docker run")
	if useRm {
		bf.WriteString(" --rm")
	}
	if dockerName != "" {
		bf.WriteString(" --name ")
		bf.WriteString(dockerName)
	}
	if network != "" {
		bf.WriteString(" --network=")
		bf.WriteString(network)
	}

	if args != nil {
		for _, arg := range args {
			bf.WriteString(" " + arg)
		}
	}

	bf.WriteString(" ")
	bf.WriteString(dockerImage)

	return exec.RunShellContext(ctx, bf.String(), progress, env.ToMap(), apps.DistRoot, true) == 0
}

func (dockerDevice) Execute(dockerName string, execCmd string, env map[string]string, progress chan string, ctx context.Context) bool {
	bf := bytes.Buffer{}
	bf.WriteString("docker exec ") // docker exec FOPS-Build-hub-fsgit-cc-fops-130 echo aaa
	for k, v := range env {
		bf.WriteString(fmt.Sprintf("-e %s=%s ", k, v))
	}
	bf.WriteString(dockerName)
	bf.WriteString(" ")
	bf.WriteString(execCmd)
	return exec.RunShellContext(ctx, bf.String(), progress, nil, apps.DistRoot, false) == 0
}

func (device dockerDevice) Copy(dockerName string, sourceFile, destFile string, env apps.EnvVO, progress chan string, ctx context.Context) bool {
	device.Execute(dockerName, "mkdir -p "+path.Dir(destFile), nil, progress, ctx)

	bf := bytes.Buffer{}
	bf.WriteString("docker cp ") // docker cp /var/lib/fops/dist/Dockerfile FOPS-Build:/var/lib/fops/dist/Dockerfile
	bf.WriteString(sourceFile)
	bf.WriteString(" ")
	bf.WriteString(dockerName)
	bf.WriteString(":")
	bf.WriteString(destFile)
	return exec.RunShellContext(ctx, bf.String(), progress, env.ToMap(), apps.DistRoot, false) == 0
}

func (dockerDevice) ExistsDocker(dockerName string) bool {
	progress := make(chan string, 1000)
	// docker inspect fops
	var exitCode = exec.RunShell(fmt.Sprintf("docker inspect %s", dockerName), progress, nil, "", false)
	lst := collections.NewListFromChan(progress)
	if exitCode != 0 {
		if lst.Contains("[]") && lst.ContainsPrefix("Error: No such object:") {
			return false
		}
		return false
	}
	if lst.Contains("[]") && lst.ContainsPrefix("Error: No such object:") {
		return false
	}
	return lst.ContainsAny(fmt.Sprintf("\"Name\": \"/%s\",", dockerName))
}

func (dockerDevice) Kill(dockerName string) {
	exec.RunShell(fmt.Sprintf("docker kill %s", dockerName), make(chan string, 1000), nil, "", false)
}

func (dockerDevice) Remove(dockerName string) {
	exec.RunShell(fmt.Sprintf("docker rm %s", dockerName), make(chan string, 1000), nil, "", false)
}

// ClearImages 清除镜像
func (dockerDevice) ClearImages(progress chan string) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始清除镜像。"

	var exitCode = exec.RunShell(`docker rmi $(docker images -f "dangling=true" -q) && docker builder prune -f && docker system prune -f`, progress, nil, "", false)
	if exitCode != 0 {
		progress <- "清除镜像镜像失败。"
		return false
	}
	progress <- "清除镜像完成。"
	return true
}

func (dockerDevice) GetVersion() string {
	receiveOutput := make(chan string, 100)
	exec.RunShell("docker version --format '{{.Server.Version}}'", receiveOutput, nil, "", false)
	lst := collections.NewListFromChan(receiveOutput)
	re := regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	for _, s := range lst.ToArray() {
		if re.MatchString(s) {
			return s
		}
	}
	return ""
}

func (dockerDevice) Login(dockerHub string, loginName string, loginPwd string, progress chan string) bool {
	if loginName != "" && loginPwd != "" {
		// 不包含域名的，意味着是登陆docker官网，不需要额外设置登陆的URL
		if !strings.Contains(dockerHub, ".") {
			dockerHub = ""
		}
		var result = exec.RunShell("docker login "+dockerHub+" -u "+loginName+" -p "+loginPwd, progress, nil, "", true)
		if result != 0 {
			progress <- "镜像仓库登陆失败。"
			return false
		}
	}

	progress <- "镜像仓库登陆成功。"
	return true
}

func (dockerDevice) Pull(image string, progress chan string) {
	exec.RunShell(fmt.Sprintf("docker pull %s", image), progress, nil, "", true)
}

func (dockerDevice) Logs(appName string, tailCount int) collections.List[string] {
	progress := make(chan string, 1000)
	// docker service logs fops
	var exitCode = exec.RunShell(fmt.Sprintf("docker logs %s --tail %d", appName, tailCount), progress, nil, "", true)
	lst := collections.NewListFromChan(progress)
	if exitCode != 0 {
		lst.Insert(0, "获取日志失败。")
	}
	return lst
}
