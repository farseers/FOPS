package request

import "time"

type WebapiVisitsRequest struct {
	AppName    string
	VisitsNode string
	StartAt    time.Time
	EndAt      time.Time
}

func (receiver *WebapiVisitsRequest) Check() {
	if receiver.EndAt.Year() == 1 {
		receiver.EndAt = time.Now()
	}
}
