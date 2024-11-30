package request

import (
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/trace"
)

type UploadRequest struct {
	List collections.List[trace.TraceContext]
}
