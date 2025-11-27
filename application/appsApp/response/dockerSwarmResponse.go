package response

import "github.com/farseer-go/docker"

type DockerSwarmResponse struct {
	docker.TaskInstanceVO
	Log string
}
