package apps

import (
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/utils/file"
	"github.com/farseer-go/utils/str"
	"time"
)

const SavePath = "/var/lib/fops/log/"

type LogQueue struct {
	BuildId  int64
	Log      string
	LogAt    time.Time
	progress chan string
}

func NewLogQueue(buildId int64) *LogQueue {
	log := &LogQueue{
		BuildId:  buildId,
		progress: make(chan string, 1000),
	}
	go log.startPush()
	return log
}

// 开启后台日志写入
func (receiver *LogQueue) startPush() {
	logfile := receiver.GenerateFilename()
	var prevContent string
	for log := range receiver.progress {
		log = flog.ClearColor(log)
		// 如果内容与前面一样，则不记录
		if prevContent == log {
			continue
		}
		prevContent = log
		//flog.Println(log)
		// 写入日志文件
		file.AppendLine(logfile, str.ToDateTime(time.Now())+" "+log)
	}
}

// View 查看日志
func (receiver *LogQueue) View() []string {
	path := SavePath + parse.ToString(receiver.BuildId) + ".txt"
	if file.IsExists(path) {
		return file.ReadAllLines(path)
	}
	return []string{}
}

// GenerateFilename 生成文件名
func (receiver *LogQueue) GenerateFilename() string {
	logfile := SavePath + parse.ToString(receiver.BuildId) + ".txt"
	// 清除历史记录（正常不会存在，当buildId被重置时，有可能会冲突）
	if file.IsExists(logfile) {
		file.Delete(logfile)
	}

	if !file.IsExists(SavePath) {
		file.CreateDir766(SavePath)
	}
	return logfile
}

// Close 关闭文件写入
func (receiver *LogQueue) Close() {
	close(receiver.progress)
}
