package job

import (
	"fops/domain/linkTrace"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/container"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	linkTraceCom "github.com/farseer-go/linkTrace"
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
	lst := repository.ToWebApiListByVisits(lastVisitsWebApiAt, endAt)
	// 按分钟做groupby
	groupBy := make(map[time.Time][]linkTraceCom.TraceContext)
	lst.GroupBy(&groupBy, func(item linkTraceCom.TraceContext) any {
		return time.Date(item.CreateAt.Year(), time.Month(item.CreateAt.Month()), item.CreateAt.Day(), item.CreateAt.Hour(), item.CreateAt.Minute(), 0, 0, time.Local)
	})

	// 根据前缀树，开始遍历统计
	lstEO := collections.NewList[linkTrace.VisitsEO]()

	// 按分钟统计
	for createAt, arrTrace := range groupBy {
		lstTrace := collections.NewList(arrTrace...)
		m := make(map[string]struct{})
		for _, item := range arrTrace {
			urlParse, _ := url.Parse(item.WebPath)
			pathPrefix := urlParse.Scheme + "://" + urlParse.Host

			paths := strings.Split(urlParse.Path, "/")
			for i := 0; i < len(paths); i++ {
				// 最后一个
				if i == len(paths)-1 {
					pathPrefix += paths[i]
				} else {
					pathPrefix += paths[i] + "/"
				}

				if _, exists := m[pathPrefix]; !exists {
					m[pathPrefix] = struct{}{}
				}
			}
		}
		// 按访问节点来遍历
		for k := range m {
			items := lstTrace.Where(func(item linkTraceCom.TraceContext) bool {
				return strings.HasPrefix(item.WebPath, k)
			}).OrderBy(func(item linkTraceCom.TraceContext) any {
				return item.UseTs.Milliseconds()
			}).ToList()
			totalCount := items.Count()
			index95 := parse.ToInt(float64(totalCount) * 0.95)
			index99 := parse.ToInt(float64(totalCount) * 0.99)

			// 设置前缀树
			lIndex := 0
			// 目录
			if strings.HasSuffix(k, "/") {
				lIndex = strings.LastIndex(k[:len(k)-1], "/")

			} else {
				lIndex = strings.LastIndex(k, "/")
			}
			visitsNodePrefix := k[:lIndex+1]
			if strings.HasSuffix(visitsNodePrefix, "//") {
				visitsNodePrefix = ""
			}
			lstEO.Add(linkTrace.VisitsEO{
				AppName:          items.First().AppName,
				CreateAt:         createAt,
				VisitsNodePrefix: visitsNodePrefix,
				VisitsNode:       k,
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

	lstEO = lstEO.OrderBy(func(item linkTrace.VisitsEO) any {
		return item.CreateAt.UnixMilli()
	}).ToList()
	_, err := repository.SaveVisitsWebApi(lstEO)
	flog.ErrorIfExists(err)
	if err == nil {
		lastVisitsWebApiAt = endAt
	}
}
