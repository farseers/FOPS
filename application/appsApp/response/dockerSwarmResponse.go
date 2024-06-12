package response

import "fops/domain/apps"

type DockerSwarmResponse struct {
	apps.DockerInstanceVO
	Log string
}
