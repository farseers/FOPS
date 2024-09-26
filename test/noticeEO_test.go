package test

import (
	"fops/domain/enum/noticeType"
	"fops/domain/monitor"
	"github.com/farseer-go/fs"
	"github.com/farseer-go/fs/modules"
	"testing"
)

func TestNoticeEO_Notice(t *testing.T) {
	fs.Initialize[modules.FarseerKernelModule]("fops-test")
	(&monitor.NoticeEO{
		NoticeType: noticeType.Telegram,
		Phone:      "",
		ApiKey:     "",
	}).Notice("测试发消息")

	(&monitor.NoticeEO{NoticeType: noticeType.Log}).Notice("测试发消息")
}
