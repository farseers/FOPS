package apps

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
	"github.com/farseer-go/fs/flog"
)

// key = nodeIP
// value = container app
var NodeDockerStatsList = collections.NewDictionary[string, collections.List[docker.DockerStatsVO]]()

type DockerInspectVO struct {
	docker.DockerStatsVO
	NodeID      string // 节点ID
	NodeName    string // 节点名称
	NodeIP      string // 集群节点
	ContainerIP string // 容器IP
	CreatedAt   string
	UpdatedAt   string
	State       string
}

// GetDockerStats 根据集群节点IP，找到对应的容器ID
func GetDockerStats(nodeIP, taskId, appName string) docker.DockerStatsVO {
	lstDockerStats := NodeDockerStatsList.GetValue(nodeIP)
	dockerStatsVO := lstDockerStats.Find(func(item *docker.DockerStatsVO) bool {
		return item.TaskId == taskId
	})

	if dockerStatsVO != nil {
		return *dockerStatsVO
	}

	//json, _ := snc.Marshal(lstDockerStats)
	flog.Infof("未找到对应的容器资源信息: %s, %s, %s", nodeIP, taskId, appName)

	return docker.DockerStatsVO{TaskId: taskId}
}
