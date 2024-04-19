// @area /linkTrace/
package linkTraceApp

import (
	"fops/application/linkTraceApp/request"
	"fops/domain/linkTrace"
	"github.com/farseer-go/collections"
)

type name struct {
}

// WebapiVisits WebApi访问统计
// @get webapiVisits
func WebapiVisits(request request.WebapiVisitsRequest, linkTraceRepository linkTrace.Repository) collections.List[linkTrace.WebapiVisitsEO] {
	return linkTraceRepository.ToWebApiVisitsList(request.AppName, request.VisitsNode, request.StartAt, request.EndAt)
}
