package request

import (
	"github.com/farseer-go/fs/trace"
)

type UploadRequest struct {
	List []trace.TraceContext
}
