package job

import (
	"fops/domain/linkTrace"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	linkTraceCom "github.com/farseer-go/linkTrace"
	"github.com/farseer-go/linkTrace/eumTraceType"
	"github.com/farseer-go/tasks"
	"net/url"
	"strings"
	"time"
)

var lastVisitsWebApiAt time.Time

// StatVisitsWebapiJob 统计webapi访问
func StatVisitsWebapiJob(*tasks.TaskContext) {
	// 缓存起来，不用每次执行都获取一次（仅在启动时获取）
	repository := container.Resolve[linkTrace.Repository]()
	if lastVisitsWebApiAt.Year() == 1 {
		lastVisitsWebApiAt, _ = repository.GetLastVisitsWebApiAt()
	}

	// 还是为1，说明从来没有执行过统计，则默认时间为昨天
	if lastVisitsWebApiAt.Year() == 1 {
		lastVisitsWebApiAt = time.Now()
		lastVisitsWebApiAt = time.Date(lastVisitsWebApiAt.Year(), lastVisitsWebApiAt.Month(), lastVisitsWebApiAt.Day(), lastVisitsWebApiAt.Hour(), lastVisitsWebApiAt.Minute(), 0, 0, time.Local)
		// 抹去秒（只要分钟）
		lastVisitsWebApiAt = lastVisitsWebApiAt.Add(-24 * time.Hour)
	}

	// 截止到当前时间的0秒
	endAt := time.Now()
	endAt = time.Date(endAt.Year(), endAt.Month(), endAt.Day(), endAt.Hour(), endAt.Minute(), 0, 0, time.Local)

	// 获取webapi链路集合
	lst := repository.ToTraceListByVisits(lastVisitsWebApiAt, endAt)
	// 按链路类型分组
	var traceTypeGroupBy map[int][]linkTraceCom.TraceContext
	lst.GroupBy(&traceTypeGroupBy, func(item linkTraceCom.TraceContext) any {
		return int(item.TraceType)
	})

	// 根据前缀树，开始遍历统计
	lstEO := collections.NewList[linkTrace.VisitsEO]()

	// 先按链路类型分组遍历
	for traceType, traceContexts := range traceTypeGroupBy {
		// 按分钟做groupBy
		groupBy := make(map[time.Time][]linkTraceCom.TraceContext)
		collections.NewList(traceContexts...).GroupBy(&groupBy, func(item linkTraceCom.TraceContext) any {
			return time.Date(item.CreateAt.Year(), time.Month(item.CreateAt.Month()), item.CreateAt.Day(), item.CreateAt.Hour(), item.CreateAt.Minute(), 0, 0, time.Local)
		})

		// 得到前缀
		for createAt, arrTrace := range groupBy {
			// 统计有多少个前缀，Value约定：0 = prefixName，1 = 该前缀的集合
			var mPathPrefix = make(map[string][]any)
			lstTrace := collections.NewList(arrTrace...)

			switch eumTraceType.Enum(traceType) {
			case eumTraceType.WebApi:
				mPathPrefix = getWebapiPrefix(lstTrace)
			case eumTraceType.MqConsumer:
				mPathPrefix["MQ"] = []any{"", lstTrace}
				for _, traceContext := range arrTrace {
					if _, exists := mPathPrefix[traceContext.ConsumerServer]; !exists {
						mPathPrefix[traceContext.ConsumerServer] = []any{"MQ", lstTrace.Where(func(item linkTraceCom.TraceContext) bool {
							return item.ConsumerServer == traceContext.ConsumerServer
						}).ToList()}
					}
					if _, exists := mPathPrefix[traceContext.ConsumerQueueName]; !exists {
						mPathPrefix[traceContext.ConsumerQueueName] = []any{traceContext.ConsumerServer, lstTrace.Where(func(item linkTraceCom.TraceContext) bool {
							return item.ConsumerQueueName == traceContext.ConsumerQueueName && item.ConsumerServer == traceContext.ConsumerServer
						}).ToList()}
					}
				}
			case eumTraceType.QueueConsumer:
				mPathPrefix["Queue"] = []any{"", lstTrace}
				for _, traceContext := range arrTrace {
					if _, exists := mPathPrefix[traceContext.ConsumerQueueName]; !exists {
						mPathPrefix[traceContext.ConsumerQueueName] = []any{"Queue", lstTrace.Where(func(item linkTraceCom.TraceContext) bool {
							return item.ConsumerQueueName == traceContext.ConsumerQueueName
						}).ToList()}
					}
				}
			case eumTraceType.FSchedule:
				mPathPrefix["FSchedule"] = []any{"", lstTrace}
				for _, traceContext := range arrTrace {
					if _, exists := mPathPrefix[traceContext.TaskGroupName]; !exists {
						mPathPrefix[traceContext.TaskGroupName] = []any{"FSchedule", lstTrace.Where(func(item linkTraceCom.TraceContext) bool {
							return item.TaskGroupName == traceContext.TaskGroupName
						}).ToList()}
					}
				}
			case eumTraceType.Task:
				mPathPrefix["Task"] = []any{"", lstTrace}
				for _, traceContext := range arrTrace {
					if _, exists := mPathPrefix[traceContext.TaskName]; !exists {
						mPathPrefix[traceContext.TaskName] = []any{"Task", lstTrace.Where(func(item linkTraceCom.TraceContext) bool {
							return item.TaskName == traceContext.TaskName
						}).ToList()}
					}
				}
			case eumTraceType.WatchKey:
				mPathPrefix["Etcd"] = []any{"", lstTrace}
				for _, traceContext := range arrTrace {
					if _, exists := mPathPrefix[traceContext.WatchKey]; !exists {
						mPathPrefix[traceContext.WatchKey] = []any{"Etcd", lstTrace.Where(func(item linkTraceCom.TraceContext) bool {
							return item.WatchKey == traceContext.WatchKey
						}).ToList()}
					}
				}
			case eumTraceType.EventConsumer:
				mPathPrefix["Event"] = []any{"", lstTrace}
				for _, traceContext := range arrTrace {
					if _, exists := mPathPrefix[traceContext.ConsumerQueueName]; !exists {
						mPathPrefix[traceContext.ConsumerQueueName] = []any{"Event", lstTrace.Where(func(item linkTraceCom.TraceContext) bool {
							return item.ConsumerQueueName == traceContext.ConsumerQueueName
						}).ToList()}
					}
				}
			}

			// 按访问节点来遍历
			for visitsNode, v := range mPathPrefix {
				visitsNodePrefix := v[0].(string)
				items := v[1].(collections.List[linkTraceCom.TraceContext])

				switch eumTraceType.Enum(traceType) {
				case eumTraceType.WebApi:
					items = lstTrace.Where(func(item linkTraceCom.TraceContext) bool {
						return strings.HasPrefix(item.WebPath, visitsNode)
					}).ToList()
				}

				items.OrderBy(func(item linkTraceCom.TraceContext) any {
					return item.UseTs.Milliseconds()
				}).ToList()

				totalCount := items.Count()
				index95 := parse.ToInt(float64(totalCount) * 0.95)
				index99 := parse.ToInt(float64(totalCount) * 0.99)

				lstEO.Add(linkTrace.VisitsEO{
					TraceType:        eumTraceType.Enum(traceType),
					AppName:          items.First().AppName,
					CreateAt:         createAt,
					VisitsNodePrefix: visitsNodePrefix,
					VisitsNode:       visitsNode,
					MinMs:            float64(items.Min(func(item linkTraceCom.TraceContext) any { return item.UseTs.Microseconds() }).(int64)) / 1000,
					MaxMs:            float64(items.Max(func(item linkTraceCom.TraceContext) any { return item.UseTs.Microseconds() }).(int64)) / 1000,
					AvgMs:            items.Average(func(item linkTraceCom.TraceContext) any { return item.UseTs.Milliseconds() }),
					Line95Ms:         float64(items.Index(index95).UseTs.Microseconds()) / 1000,
					Line99Ms:         float64(items.Index(index99).UseTs.Microseconds()) / 1000,
					ErrorCount:       items.Where(func(item linkTraceCom.TraceContext) bool { return item.Exception != nil && !item.Exception.IsNil() }).Count(),
					TotalCount:       totalCount,
					QPS:              float64(totalCount) / 60,
				})
			}
		}
	}

	lstEO = lstEO.OrderBy(func(item linkTrace.VisitsEO) any {
		return item.CreateAt.UnixMilli()
	}).ToList()
	_, err := repository.SaveVisitsWebApi(lstEO)
	flog.ErrorIfExists(err)
	if err == nil {
		lastVisitsWebApiAt = endAt
	}
}

func getWebapiPrefix(lstTrace collections.List[linkTraceCom.TraceContext]) map[string][]any {
	mPathPrefix := make(map[string][]any)
	for _, item := range lstTrace.ToArray() {
		urlParse, _ := url.Parse(item.WebPath)
		path := urlParse.Scheme + "://" + urlParse.Host

		paths := strings.Split(urlParse.Path, "/")
		for i := 0; i < len(paths); i++ {
			// 最后一个
			if i == len(paths)-1 {
				path += paths[i]
			} else {
				path += paths[i] + "/"
			}

			if _, exists := mPathPrefix[path]; !exists {
				// 设置前缀树
				lIndex := 0
				// 目录
				if strings.HasSuffix(path, "/") {
					lIndex = strings.LastIndex(path[:len(path)-1], "/")

				} else {
					lIndex = strings.LastIndex(path, "/")
				}
				visitsNodePrefix := path[:lIndex+1]
				if strings.HasSuffix(visitsNodePrefix, "//") {
					visitsNodePrefix = ""
				}
				mPathPrefix[path] = []any{visitsNodePrefix, lstTrace.Where(func(item linkTraceCom.TraceContext) bool {
					return strings.HasPrefix(item.WebPath, path)
				}).ToList()}
			}
		}
	}
	return mPathPrefix
}
