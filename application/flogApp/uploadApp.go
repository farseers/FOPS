// @area /flog/
package flogApp

import (
	"fops/application/flogApp/request"
	"fops/domain/logData"
	"github.com/farseer-go/fs/exception"
)

// Upload 上传链路记录
// @post upload
func Upload(req request.UploadRequest, logDataRepository logData.Repository) {
	err := logDataRepository.Save(req.List)
	exception.ThrowWebExceptionError(403, err)
}
