// @area /flog/
package flogApp

import (
	"fops/application/flogApp/request"
	"fops/domain/logData"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/queue"
)

// Upload 上传链路记录
// @post upload
func Upload(req request.UploadRequest, logDataRepository logData.Repository) {
	// 写入到本地队列
	req.List.Foreach(func(item *flog.LogData) {
		queue.Push("flog", *item)
	})
}
