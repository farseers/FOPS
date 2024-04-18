// @area /linkTrace/
package linkTraceApp

import (
	"fops/application/linkTraceApp/request"
	"fops/domain/linkTrace"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/parse"
	linkTraceCom "github.com/farseer-go/linkTrace"
	"net/url"
	"strings"
)

type name struct {
}

// WebapiVisits WebApi访问统计
// @get webapiVisits
func WebapiVisits(request request.WebapiVisitsRequest, linkTraceRepository linkTrace.Repository) collections.List[linkTrace.WebapiVisitsEO] {
	m := make(map[string]struct{})
	lst := linkTraceRepository.ToWebApiVisitsList(request.AppName, request.VisitsNode, request.StartAt, request.EndAt)
	// 先得到所有的前缀树
	lst.Foreach(func(item *linkTraceCom.TraceContext) {
		urlParse, _ := url.Parse(item.WebPath)
		pathPrefix := urlParse.Scheme + "://" + urlParse.Host

		// 只统计根目录
		if request.VisitsNode == "" {
			pathPrefix += "/"
			if _, exists := m[pathPrefix]; !exists {
				m[pathPrefix] = struct{}{}
			}
		} else {
			paths := strings.Split(urlParse.Path, "/")
			for i := 0; i < len(paths); i++ {
				// 最后一个
				if i == len(paths)-1 {
					pathPrefix += paths[i]
				} else {
					pathPrefix += paths[i] + "/"
				}

				if request.VisitsNode == pathPrefix {
					if _, exists := m[pathPrefix]; !exists {
						m[pathPrefix] = struct{}{}
					}
				}
			}
		}
	})

	// 根据前缀树，开始遍历统计
	lstEO := collections.NewList[linkTrace.WebapiVisitsEO]()
	for k := range m {
		items := lst.Where(func(item linkTraceCom.TraceContext) bool {
			return strings.HasPrefix(item.WebPath, k)
		}).OrderBy(func(item linkTraceCom.TraceContext) any {
			return item.UseTs.Milliseconds()
		}).ToList()
		totalCount := items.Count()
		index95 := parse.ToInt(float64(totalCount) * 0.95)
		index99 := parse.ToInt(float64(totalCount) * 0.99)

		lstEO.Add(linkTrace.WebapiVisitsEO{
			VisitsNode:    k,
			SpeedMinMs:    items.Min(func(item linkTraceCom.TraceContext) any { return item.UseTs.Milliseconds() }).(int64),
			SpeedMaxMs:    items.Max(func(item linkTraceCom.TraceContext) any { return item.UseTs.Milliseconds() }).(int64),
			SpeedAvgMs:    items.Average(func(item linkTraceCom.TraceContext) any { return item.UseTs.Milliseconds() }),
			Speed95LineMs: items.Index(index95).UseTs.Milliseconds(),
			Speed99LineMs: items.Index(index99).UseTs.Milliseconds(),
			ErrorCount:    items.Where(func(item linkTraceCom.TraceContext) bool { return item.Exception != nil && !item.Exception.IsNil() }).Count(),
			TotalCount:    totalCount,
			QPS:           float64(totalCount) / request.EndAt.Sub(request.StartAt).Seconds(),
		})
	}
	return lstEO.OrderBy(func(item linkTrace.WebapiVisitsEO) any {
		return item.VisitsNode
	}).ToList()
}
