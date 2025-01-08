package http

import (
	"fops/domain/fSchedule"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/utils/http"
)

type FScheduleHttp struct {
}

func (receiver FScheduleHttp) StatList(addr string) collections.List[fSchedule.StatTaskEO] {
	url := addr + "/basicapi/stat/statList"
	apiResponse, err := http.GetJson[core.ApiResponse[collections.List[fSchedule.StatTaskEO]]](url, nil, 500)
	if err != nil {
		flog.Warningf("请求接口：%s 失败，错误信息：%s", url, err.Error())
	}
	return apiResponse.Data
}
