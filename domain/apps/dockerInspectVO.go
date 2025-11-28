package apps

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/docker"
)

// key = nodeIP
// value = container app
var NodeDockerStatsList = collections.NewDictionary[string, collections.List[docker.DockerStatsVO]]()

type DockerInspectVO struct {
	docker.DockerStatsVO
	TaskId      string // 任务ID（docker service ps xxx 得到）
	Node        string // 节点
	NodeIP      string // 集群节点
	ContainerIP string // 容器IP
	CreatedAt   string
	UpdatedAt   string
	State       string
}

// GetDockerStats 根据集群节点IP，找到对应的容器ID
func GetDockerStats(nodeIP, taskId string) docker.DockerStatsVO {
	lstDockerStats := NodeDockerStatsList.GetValue(nodeIP)
	dockerStatsVO := lstDockerStats.Find(func(item *docker.DockerStatsVO) bool {
		return item.TaskId == taskId
	})

	if dockerStatsVO != nil {
		return *dockerStatsVO
	}
	return docker.DockerStatsVO{TaskId: taskId}
}
