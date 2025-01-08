// @area /flog/
package flogApp

import (
	"fops/application/flogApp/request"
	"fops/domain/logData"

	"github.com/farseer-go/fs/trace"
	"github.com/farseer-go/queue"
)

// Upload 上传链路记录
// @post upload
func Upload(req request.UploadRequest, logDataRepository logData.Repository) {
	if t := trace.CurTraceContext.Get(); t != nil {
		t.Ignore()
	}

	// 先发送到本地队列
	for _, item := range req.List {
		queue.Push("flog", &item)
	}

	req.List = nil
}
