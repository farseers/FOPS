// @area /linkTrace/
package linkTraceApp

import (
	"fops/application/linkTraceApp/request"
	"fops/domain/linkTrace"
	"github.com/farseer-go/collections"
)

type name struct {
}

// Visits 访问统计
// @get visits
func Visits(request request.VisitsRequest, linkTraceRepository linkTrace.Repository) collections.List[linkTrace.VisitsEO] {
	return linkTraceRepository.ToVisitsList(request.AppName, request.VisitsNode, request.StartAt, request.EndAt)
}
