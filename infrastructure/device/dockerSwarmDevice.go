package device

import (
	"context"
	"fmt"
	"fops/domain/apps"
	"fops/domain/cluster"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/utils/exec"
	"strings"
)

func RegisterDockerSwarmDevice() {
	container.Register(func() apps.IDockerSwarmDevice { return &dockerSwarmDevice{} })
}

type dockerSwarmDevice struct {
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

func (dockerSwarmDevice) SetImagesAndReplicas(cluster cluster.DomainObject, appName string, dockerImages string, dockerReplicas int, progress chan string) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始更新Docker Swarm的镜像版本。"

	var exitCode = exec.RunShell(fmt.Sprintf("docker service update --image %s --replicas %v --update-delay 10s --with-registry-auth %s", dockerImages, dockerReplicas, appName), progress, nil, "", false)
	if exitCode != 0 {
		progress <- "Docker Swarm更新镜像失败。"
		return false
	}
	progress <- "Docker Swarm更新镜像版本完成。"
	return true
}

func (dockerSwarmDevice) SetImages(cluster cluster.DomainObject, appName string, dockerImages string, progress chan string) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始更新Docker Swarm的镜像版本。"

	var exitCode = exec.RunShell(fmt.Sprintf("docker service update --image %s --update-delay 10s --with-registry-auth %s", dockerImages, appName), progress, nil, "", false)
	if exitCode != 0 {
		progress <- "Docker Swarm更新镜像失败。"
		return false
	}
	progress <- "Docker Swarm更新镜像版本完成。"
	return true
}

func (dockerSwarmDevice) SetReplicas(cluster cluster.DomainObject, appName string, dockerReplicas int, progress chan string) bool {
	progress <- "---------------------------------------------------------"
	progress <- "开始更新Docker Swarm的副本数量。"

	var exitCode = exec.RunShell(fmt.Sprintf("docker service update --replicas %v --with-registry-auth %s", dockerReplicas, appName), progress, nil, "", false)
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

	var exitCode = exec.RunShell(fmt.Sprintf("docker service update --with-registry-auth --force %s", appName), progress, nil, "", false)
	if exitCode != 0 {
		progress <- "Docker Swarm的容器重启失败。"
		return false
	}
	progress <- "Docker Swarm的容器重启完成。"
	return true
}

func (dockerSwarmDevice) ExistsDocker(appName string) bool {
	progress := make(chan string, 1000)
	// docker service inspect fops
	var exitCode = exec.RunShell(fmt.Sprintf("docker service inspect %s", appName), progress, nil, "", true)
	lst := collections.NewListFromChan(progress)
	if exitCode != 0 {
		if lst.Contains("[]") && lst.ContainsPrefix("Status: Error: no such service:") {
			return false
		}
		progress <- "获取应用信息时失败。"
		return false
	}
	if lst.Contains("[]") && lst.ContainsPrefix("Status: Error: no such service:") {
		return false
	}
	return lst.ContainsAny(fmt.Sprintf("\"Name\": \"%s\"", appName))
}

func (dockerSwarmDevice) CreateService(appName, dockerNodeRole, additionalScripts, dockerNetwork string, dockerReplicas int, dockerImages string, progress chan string, ctx context.Context) bool {
	progress <- "开始创建Docker Swarm容器服务。"

	shell := fmt.Sprintf("docker service create --with-registry-auth --name %s --replicas %v -d --network=%s --constraint node.role==%s --mount type=bind,src=/etc/localtime,dst=/etc/localtime %s %s", appName, dockerReplicas, dockerNetwork, dockerNodeRole, additionalScripts, dockerImages)
	var exitCode = exec.RunShellContext(ctx, shell, progress, nil, "", true)
	if exitCode != 0 {
		progress <- "创建Docker Swarm容器失败了。"
		return false
	}
	return true
}

func (dockerSwarmDevice) Logs(appName string, tailCount int) collections.List[string] {
	progress := make(chan string, 1000)
	// docker service logs fops
	var exitCode = exec.RunShell(fmt.Sprintf("docker service logs %s --tail %d", appName, tailCount), progress, nil, "", true)
	lst := collections.NewListFromChan(progress)
	if exitCode != 0 {
		lst.Insert(0, "获取日志失败。")
	}
	return lst
}

func (dockerSwarmDevice) ServiceList() collections.List[apps.DockerServiceVO] {
	progress := make(chan string, 1000)
	// docker service ls
	var exitCode = exec.RunShell("docker service ls", progress, nil, "", false)
	serviceList := collections.NewListFromChan(progress)
	lstDockerName := collections.NewList[apps.DockerServiceVO]()
	if exitCode != 0 || serviceList.Count() == 0 {
		return lstDockerName
	}

	// 移除标题
	serviceList.RemoveAt(0)
	serviceList.Foreach(func(service *string) {
		// 移除容器ID
		*service = strings.TrimSpace((*service)[12:])
		// 移除空格
		*service = strings.Replace(*service, "\t", "", -1)
		sers := collections.NewList(strings.Split(*service, " ")...)
		sers.RemoveAll(func(item string) bool {
			return item == ""
		})

		// redis|replicated|1/1|redis:latest
		// 满足长度格式才继续
		if sers.Count() != 4 {
			return
		}
		lstDockerName.Add(apps.DockerServiceVO{
			Name:      sers.Index(0),
			Instances: parse.ToInt(strings.Split(sers.Index(2), "/")[0]),
			Replicas:  parse.ToInt(strings.Split(sers.Index(2), "/")[1]),
			Image:     sers.Index(3),
		})
	})
	return lstDockerName
}

func (dockerSwarmDevice) PS(appName string) collections.List[apps.DockerInstanceVO] {
	progress := make(chan string, 1000)
	// docker service ps fops
	var exitCode = exec.RunShell(fmt.Sprintf("docker service ps %s", appName), progress, nil, "", false)
	serviceList := collections.NewListFromChan(progress)
	lstDockerInstance := collections.NewList[apps.DockerInstanceVO]()
	if exitCode != 0 || serviceList.Count() == 0 {
		return lstDockerInstance
	}

	// 移除标题
	serviceList.RemoveAt(0)
	serviceList.Foreach(func(service *string) {
		// 移除空格
		*service = strings.Replace(*service, "\t", "", -1)
		sers := collections.NewList(strings.Split(*service, " ")...)
		sers.RemoveAll(func(item string) bool {
			return item == ""
		})

		// k0d7jnwrr8st|fops.1|hub.fsgit.cc/hub:fops.551|test|Running|Running 4 minutes ago
		// 满足长度格式才继续
		if sers.Count() < 6 {
			return
		}
		vo := apps.DockerInstanceVO{
			Id:        sers.Index(0),
			Name:      sers.Index(1),
			Image:     sers.Index(2),
			Node:      sers.Index(3),
			State:     sers.Index(4),
			StateInfo: sers.Index(5),
		}
		if sers.Count() > 6 {
			vo.Error = sers.Index(6)
		}
		lstDockerInstance.Add(vo)
	})
	return lstDockerInstance
}
