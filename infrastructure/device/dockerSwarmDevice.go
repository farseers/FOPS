package device

import (
	"context"
	"fmt"
	"fops/domain/apps"
	"fops/domain/cluster"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/utils/exec"
)

func RegisterDockerSwarmDevice() {
	container.Register(func() apps.IDockerSwarmDevice { return &dockerSwarmDevice{} })
}

type dockerSwarmDevice struct {
}

func (dockerSwarmDevice) SetImages(cluster cluster.DomainObject, appName string, dockerImages string, progress chan string, ctx context.Context) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始更新Docker Swarm的镜像版本。"

	var exitCode = exec.RunShellContext(ctx, fmt.Sprintf("docker service update --image %s %s", dockerImages, appName), progress, nil, "", false)
	if exitCode != 0 {
		progress <- "Docker Swarm更新镜像失败。"
		return false
	}
	progress <- "Docker Swarm更新镜像版本完成。"
	return true
}

func (dockerSwarmDevice) ExistsDocker(cluster cluster.DomainObject, appName string) bool {
	progress := make(chan string, 1000)
	// docker service inspect fops
	var exitCode = exec.RunShell(fmt.Sprintf("docker service inspect %s", appName), progress, nil, "", false)
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

func (dockerSwarmDevice) CreateService(appName, dockerNodeRole, additionalScripts, dockerNetwork string, dockerReplicas int, dockerImages string, progress chan string, ctx context.Context) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始创建Docker Swarm容器服务。"

	shell := fmt.Sprintf("docker service create --name %s --replicas %v -d --network=%s --constraint node.role==%s --mount type=bind,src=/etc/localtime,dst=/etc/localtime %s %s", appName, dockerReplicas, dockerNetwork, dockerNodeRole, additionalScripts, dockerImages)
	var exitCode = exec.RunShellContext(ctx, shell, progress, nil, "", false)
	if exitCode != 0 {
		progress <- "创建Docker Swarm容器失败了。"
		return false
	}
	return true
}

func (dockerSwarmDevice) DeleteService(appName string, progress chan string) bool {
	// docker service rm fops
	var exitCode = exec.RunShell(fmt.Sprintf("docker service rm %s", appName), progress, nil, "", false)
	if exitCode != 0 {
		progress <- "删除Docker Swarm容器失败了。"
		return false
	}
	return true
}

func (dockerSwarmDevice) SetReplicas(cluster cluster.DomainObject, appName string, dockerReplicas int, progress chan string) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始更新Docker Swarm的副本数量。"

	var exitCode = exec.RunShell(fmt.Sprintf("docker service update --replicas %v %s", dockerReplicas, appName), progress, nil, "", false)
	if exitCode != 0 {
		progress <- "Docker Swarm的副本数量更新失败。"
		return false
	}
	progress <- "Docker Swarm的副本数量更新完成。"
	return true
}

func (dockerSwarmDevice) Restart(cluster cluster.DomainObject, appName string, progress chan string) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始重启Docker Swarm的容器。"

	var exitCode = exec.RunShell(fmt.Sprintf("docker service update --force %s", appName), progress, nil, "", false)
	if exitCode != 0 {
		progress <- "Docker Swarm的容器重启失败。"
		return false
	}
	progress <- "Docker Swarm的容器重启完成。"
	return true
}
