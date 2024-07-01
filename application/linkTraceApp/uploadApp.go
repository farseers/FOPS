// @area /linkTrace/
package linkTraceApp

import (
	"fops/application/linkTraceApp/request"
	"fops/domain/linkTrace"
	"github.com/farseer-go/fs/trace"
	linkTraceCom "github.com/farseer-go/linkTrace"
	"github.com/farseer-go/queue"
)

// Upload 上传链路记录
// @post upload
func Upload(req request.UploadRequest, linkTraceRepository linkTrace.Repository) {
	if t := trace.CurTraceContext.Get(); t != nil {
		t.Ignore()
	}

	return
	// 先发送到本地队列
	req.List.Foreach(func(item *linkTraceCom.TraceContext) {
		queue.Push("linkTrace", item)
	})

	req.List.Clear()
}
