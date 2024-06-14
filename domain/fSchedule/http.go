package fSchedule

import "github.com/farseer-go/collections"

type Http interface {
	// StatList 任务执行统计
	StatList(addr string) collections.List[StatTaskEO]
}

type StatTaskEO struct {
	ClientName    string // 应用名称
	ExecuteStatus int    // 状态
	Count         int    // 日志数量
}
