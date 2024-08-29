package device

import (
	"bytes"
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

func (dockerSwarmDevice) Restart(appName string, progress chan string) bool {
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

func (dockerSwarmDevice) CreateService(appName, dockerNodeRole, additionalScripts, dockerNetwork string, dockerReplicas int, dockerImages string, limitCpus float64, limitMemory string, progress chan string, ctx context.Context) bool {
	progress <- "开始创建Docker Swarm容器服务。"

	var sb bytes.Buffer
	sb.WriteString("docker service create --with-registry-auth --mount type=bind,src=/etc/localtime,dst=/etc/localtime")
	sb.WriteString(fmt.Sprintf(" --name %s --replicas %v -d --network=%s", appName, dockerReplicas, dockerNetwork))

	// 所有节点都要运行
	if dockerNodeRole == "global" {
		sb.WriteString(" --mode global")
	} else {
		sb.WriteString(fmt.Sprintf(" --constraint node.role==%s", dockerNodeRole))
	}

	if limitCpus > 0 {
		sb.WriteString(fmt.Sprintf(" --limit-cpu=%f", limitCpus))
	}
	if limitMemory != "" {
		sb.WriteString(fmt.Sprintf(" --limit-memory=%s", limitMemory))
	}
	sb.WriteString(fmt.Sprintf(" %s %s", additionalScripts, dockerImages))

	var exitCode = exec.RunShellContext(ctx, sb.String(), progress, nil, "", true)
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
	lst.Foreach(func(item *string) {
		if strings.Contains(*item, "|") {
			*item = strings.TrimSpace(strings.SplitN(*item, "|", 2)[1])
		}
	})
	return lst
}

func (dockerSwarmDevice) ServiceList() collections.List[apps.DockerServiceVO] {
	progress := make(chan string, 1000)
	// docker service ls --format "table {{.ID}}|{{.Name}}|{{.Mode}}|{{.Replicas}}|{{.Image}}|{{.Ports}}"
	var exitCode = exec.RunShell("docker service ls --format \"table {{.ID}}|{{.Name}}|{{.Replicas}}|{{.Image}}\"", progress, nil, "", false)
	serviceList := collections.NewListFromChan(progress)
	lstDockerName := collections.NewList[apps.DockerServiceVO]()
	if exitCode != 0 || serviceList.Count() == 0 {
		return lstDockerName
	}

	// 移除标题
	serviceList.RemoveAt(0)
	serviceList.Foreach(func(service *string) {
		// vwceboa7gtmu|redis|1/1|redis:latest
		sers := strings.Split(*service, "|")
		if len(sers) < 4 {
			return
		}
		lstDockerName.Add(apps.DockerServiceVO{
			Id:        sers[0],
			Name:      sers[1],
			Instances: parse.ToInt(strings.Split(sers[2], "/")[0]),
			Replicas:  parse.ToInt(strings.Split(sers[2], "/")[1]),
			Image:     sers[3],
		})
	})
	return lstDockerName
}

func (dockerSwarmDevice) PS(appName string) collections.List[apps.DockerInstanceVO] {
	progress := make(chan string, 1000)
	// docker service ps fops --format "table {{.ID}}|{{.Name}}|{{.Image}}|{{.Node}}|{{.DesiredState}}|{{.CurrentState}}|{{.Error}}"
	var exitCode = exec.RunShell(fmt.Sprintf("docker service ps %s --format \"table {{.ID}}|{{.Name}}|{{.Image}}|{{.Node}}|{{.DesiredState}}|{{.CurrentState}}|{{.Error}}\"", appName), progress, nil, "", false)
	serviceList := collections.NewListFromChan(progress)
	lstDockerInstance := collections.NewList[apps.DockerInstanceVO]()
	if exitCode != 0 || serviceList.Count() == 0 {
		return lstDockerInstance
	}

	// 移除标题
	serviceList.RemoveAt(0)
	serviceList.Foreach(func(service *string) {
		// whw9erkpysrj|fops|fops.552|test|Running|Running 17 minutes ago|
		sers := strings.Split(*service, "|")
		if len(sers) < 7 {
			return
		}
		lstDockerInstance.Add(apps.DockerInstanceVO{
			Id:        sers[0],
			Name:      sers[1],
			Image:     sers[2],
			Node:      sers[3],
			State:     sers[4],
			StateInfo: sers[5],
			Error:     sers[6],
		})
	})
	return lstDockerInstance
}

func (dockerSwarmDevice) NodeList() collections.List[apps.DockerNodeVO] {
	progress := make(chan string, 1000)
	// docker node ls --format "table {{.Hostname}}|{{.Status}}|{{.Availability}}|{{.ManagerStatus}}|{{.EngineVersion}}"
	var exitCode = exec.RunShell("docker node ls --format \"table {{.Hostname}}|{{.Status}}|{{.Availability}}|{{.ManagerStatus}}|{{.EngineVersion}}\"", progress, nil, "", false)
	serviceList := collections.NewListFromChan(progress)
	lstDockerInstance := collections.NewList[apps.DockerNodeVO]()
	if exitCode != 0 || serviceList.Count() == 0 {
		return lstDockerInstance
	}

	// 移除标题
	serviceList.RemoveAt(0)
	serviceList.Foreach(func(service *string) {
		// test|Ready|Active|Leader|20.10.17
		sers := strings.Split(*service, "|")
		if len(sers) < 5 {
			return
		}
		lstDockerInstance.Add(apps.DockerNodeVO{
			NodeName:      sers[0],
			Status:        sers[1],
			Availability:  sers[2],
			IsMaster:      sers[3] == "Leader",
			EngineVersion: sers[4],
		})
	})
	return lstDockerInstance
}

func (dockerSwarmDevice) NodeInfo(nodeName string) apps.DockerNodeVO {
	progress := make(chan string, 1000)
	// docker node inspect node_1 --pretty
	var exitCode = exec.RunShell(fmt.Sprintf("docker node inspect %s --pretty", nodeName), progress, nil, "", false)
	serviceList := collections.NewListFromChan(progress)
	vo := apps.DockerNodeVO{
		Label: collections.NewList[apps.DockerLabelVO](),
	}
	if exitCode != 0 || serviceList.Count() == 0 {
		return vo
	}
	serviceList.For(func(index int, item *string) {
		kv := strings.Split(*item, ":")
		if len(kv) != 2 {
			return
		}
		name := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		switch name {
		case "Address":
			vo.IP = val
		case "Operating System":
			vo.OS = val
		case "Architecture":
			vo.Architecture = val
		case "CPUs":
			vo.CPUs = val
		case "Memory":
			vo.Memory = val
		case "Labels":
			// 标签要特殊处理
			/*
			   Labels:
			    - run=job
			    - type=master
			*/
			tag := " - "
			for {
				index++
				content := serviceList.Index(index)
				if !strings.HasPrefix(content, tag) {
					return
				}
				// 移除标签
				content = strings.TrimSpace(content[len(tag):])
				vo.Label.Add(apps.DockerLabelVO{
					Name:  strings.Split(content, "=")[0],
					Value: strings.Split(content, "=")[1],
				})
			}
		}
	})
	return vo
}
