package response

import "github.com/farseer-go/docker"

type DockerSwarmResponse struct {
	docker.ServicePsVO
	Log string
}
