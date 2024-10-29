// @area /monitor/
package monitorApp

import (
	"fops/application/monitorApp/request"
	"fops/application/monitorApp/response"
	"fops/domain/enum/noticeType"
	"fops/domain/monitor"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
)

// DropDownListAppInfo 应用下拉框
// @post appList
// @filter application.Jwt
func DropDownListAppInfo(monitorRepository monitor.Repository) collections.List[response.AppInfoResponse] {
	lst := monitorRepository.DropDownListAppInfo()
	return mapper.ToList[response.AppInfoResponse](lst)
}

// ToListPageRule 规则分页
// @post ruleList
// @filter application.Jwt
func ToListPageRule(appName string, pageSize, pageIndex int, monitorRepository monitor.Repository) collections.PageList[response.RuleResponse] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	lst := monitorRepository.ToListPageRule(appName, pageSize, pageIndex)
	resList := mapper.ToPageList[response.RuleResponse](lst)
	resList.List.Foreach(func(item *response.RuleResponse) {
		if len(item.NoticeIds) > 0 {
			item.NoticeList = monitorRepository.ToListNoticeById(item.NoticeIds)
		}
	})
	return resList
}

// DeleteRule 删除规则
// @post delRule
// @filter application.Jwt
func DeleteRule(id int64, monitorRepository monitor.Repository) {
	err := monitorRepository.DeleteRule(id)
	exception.ThrowWebExceptionError(403, err)
}

// ToEntityRule 规则详情
// @post infoRule
// @filter application.Jwt
func ToEntityRule(id int64, monitorRepository monitor.Repository) response.RuleResponse {
	info := monitorRepository.ToEntityRule(id)
	resInfo := mapper.Single[response.RuleResponse](info)
	if len(info.NoticeIds) > 0 {
		resInfo.NoticeList = monitorRepository.ToListNoticeById(info.NoticeIds)
	}
	return resInfo
}

// SaveRule 保存规则
// @post saveRule
// @filter application.Jwt
func SaveRule(req request.SaveRuleRequest, monitorRepository monitor.Repository) {
	do := mapper.Single[monitor.RuleEO](req)
	if req.Id > 0 {
		// 更新
		err := monitorRepository.UpdateRule(req.Id, do)
		exception.ThrowWebExceptionError(403, err)
	} else {
		// 添加规则
		err := monitorRepository.AddRule(do)
		exception.ThrowWebExceptionError(403, err)
	}
}

// ToListPageNotice 通知人列表
// @post noticeList
// @filter application.Jwt
func ToListPageNotice(name string, pageSize, pageIndex int, monitorRepository monitor.Repository) collections.PageList[monitor.NoticeEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return monitorRepository.ToListPageNotice(name, pageSize, pageIndex)
}

// DeleteNotice 删除通知人
// @post delNotice
// @filter application.Jwt
func DeleteNotice(id int64, monitorRepository monitor.Repository) {
	err := monitorRepository.DeleteNotice(id)
	exception.ThrowWebExceptionError(403, err)
}

// ToEntityNotice 通知人详情
// @post infoNotice
// @filter application.Jwt
func ToEntityNotice(id int64, monitorRepository monitor.Repository) monitor.NoticeEO {
	return monitorRepository.ToEntityNotice(id)
}

// SaveNotice 保存通知人
// @post saveNotice
// @filter application.Jwt
func SaveNotice(req monitor.NoticeEO, monitorRepository monitor.Repository) {
	if req.Id > 0 {
		// 更新
		err := monitorRepository.UpdateNotice(req.Id, req)
		exception.ThrowWebExceptionError(403, err)
	} else {
		// 添加规则
		err := monitorRepository.AddRNotice(req)
		exception.ThrowWebExceptionError(403, err)
	}
}

// ToListPageData 监控数据列表
// @post dataList
// @filter application.Jwt
func ToListPageData(appName string, pageSize, pageIndex int, monitorRepository monitor.Repository) collections.PageList[monitor.DataEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return monitorRepository.ToListPageData(appName, pageSize, pageIndex)
}

// ToListPageNoticeLog 通知消息列表
// @post noticeLogList
// @filter application.Jwt
func ToListPageNoticeLog(appName string, pageSize, pageIndex int, monitorRepository monitor.Repository) collections.PageList[monitor.NoticeLogEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return monitorRepository.ToListPageNoticeLog(appName, pageSize, pageIndex)
}

// ToListPageNoticeLog 通知消息未读消息列表
// @post noticeLogNoReadList
// @filter application.Jwt
func ToListPageNoticeLogNoRead(top int, monitorRepository monitor.Repository) collections.List[monitor.NoticeLogEO] {
	return monitorRepository.ToListPageNoticeLogNoRead(top)
}

// UpdateNoticeLogRead 全部已读设置
// @post allRead
// @filter application.Jwt
func UpdateNoticeLogRead(req request.ReadRequest, monitorRepository monitor.Repository) {
	err := monitorRepository.UpdateNoticeLogRead(req.Ids)
	exception.ThrowRefuseExceptionError(err)
}

// DeleteNoticeLog 删除通知消息日志
// @post delNoticeLog
// @filter application.Jwt
func DeleteNoticeLog(monitorRepository monitor.Repository) {
	// 7天之前的全部删除
	err := monitorRepository.DeleteNoticeLog(time.Now().AddDate(0, 0, -7))
	exception.ThrowWebExceptionError(403, err)
}

// DrpBaseList drp基础类型列表
// @post drpBaseList
// @filter application.Jwt
func DrpBaseList(baseType string) map[string][]response.KeyValueResponse {

	reqList := collections.NewList(strings.Split(baseType, ",")...)
	resMap := make(map[string][]response.KeyValueResponse)
	if reqList.Contains("1") {
		lst := noticeType.ToList()
		resList := collections.NewList[response.KeyValueResponse]()
		lst.Foreach(func(item *noticeType.Enum) {
			resList.Add(response.KeyValueResponse{
				Key:   int(*item),
				Value: item.ToString(),
			})
		})
		resMap["NoticeTypeList"] = resList.ToArray()
	}
	if reqList.Contains("2") {
		resList := collections.NewList[response.KeyValueResponse]()
		resList.Add(response.KeyValueResponse{
			Key:   1,
			Value: ">",
		})
		resList.Add(response.KeyValueResponse{
			Key:   2,
			Value: "<",
		})
		resList.Add(response.KeyValueResponse{
			Key:   3,
			Value: "=",
		})
		resList.Add(response.KeyValueResponse{
			Key:   4,
			Value: "!=",
		})
		resMap["CompareList"] = resList.ToArray()
	}
	return resMap
}
