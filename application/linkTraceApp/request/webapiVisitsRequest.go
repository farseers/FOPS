package request

import "time"

type WebapiVisitsRequest struct {
	AppName    string
	VisitsNode string
	StartAt    time.Time
	EndAt      time.Time
}
