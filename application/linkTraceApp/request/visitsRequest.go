package request

import "time"

type VisitsRequest struct {
	AppName    string
	VisitsNode string
	StartAt    time.Time
	EndAt      time.Time
}

func (receiver *VisitsRequest) Check() {
	if receiver.EndAt.Year() == 1 {
		receiver.EndAt = time.Now()
	}
}
