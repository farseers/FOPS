package response

import "github.com/farseer-go/docker"

type DockerSwarmResponse struct {
	docker.ServiceTaskVO
	Log string
}
