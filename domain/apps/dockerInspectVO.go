package apps

import "time"

type DockerInspectVO struct {
	ID        string
	IP        string
	CreatedAt time.Time
	UpdatedAt time.Time
	State     string
}
