package device

import (
	"fops/domain/apps"
	"strconv"

	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/utils/str"
)

func RegisterDockerDevice() {
	container.Register(func() apps.IDockerDevice {
		return &dockerDevice{}
	})
}

type dockerDevice struct {
}

func (dockerDevice) GetDockerHub(dockerHubAddress string) string {
	var dockerHub = "localhost"
	if dockerHubAddress != "" {
		dockerHub = str.CutRight(dockerHubAddress, "/")
	}
	return dockerHub
}

func (receiver dockerDevice) GetDockerImage(dockerHubAddress string, appName string, buildNumber int) string {
	return receiver.GetDockerHub(dockerHubAddress) + ":" + appName + "." + strconv.Itoa(buildNumber)
}
