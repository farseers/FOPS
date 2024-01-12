// @area /linkTrace/
package linkTraceApp

import (
	"fops/application/linkTraceApp/request"
	"fops/domain/linkTrace"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/trace"
)

// Upload 上传链路记录
// @post upload
func Upload(req request.UploadRequest, linkTraceRepository linkTrace.Repository) {
	if t := trace.CurTraceContext.Get(); t != nil {
		t.Ignore()
	}

	err := linkTraceRepository.Save(req.List)
	exception.ThrowWebExceptionError(403, err)
}
