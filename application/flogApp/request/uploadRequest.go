package request

import (
	"github.com/farseer-go/fs/flog"
)

type UploadRequest struct {
	List []flog.LogData
}
