// @area /monitor/
package monitorApp

import (
	"fops/application/monitorApp/request"
	"fops/application/monitorApp/response"
	"fops/domain/monitor"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
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
func ToListPageRule(pageSize, pageIndex int, monitorRepository monitor.Repository) collections.PageList[monitor.RuleEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return monitorRepository.ToListPageRule(pageSize, pageIndex)
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
func ToEntityRule(id int64, monitorRepository monitor.Repository) monitor.RuleEO {
	return monitorRepository.ToEntityRule(id)
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
		do.Enable = true
		err := monitorRepository.AddRule(do)
		exception.ThrowWebExceptionError(403, err)
	}
}

// ToListPageNotice 通知人列表
// @post noticeList
// @filter application.Jwt
func ToListPageNotice(pageSize, pageIndex int, monitorRepository monitor.Repository) collections.PageList[monitor.NoticeEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return monitorRepository.ToListPageNotice(pageSize, pageIndex)
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
		req.Enable = true
		err := monitorRepository.AddRNotice(req)
		exception.ThrowWebExceptionError(403, err)
	}
}

// ToListPageData 监控数据列表
// @post dataList
// @filter application.Jwt
func ToListPageData(appId string, pageSize, pageIndex int, monitorRepository monitor.Repository) collections.PageList[monitor.DataEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return monitorRepository.ToListPageData(appId, pageSize, pageIndex)
}

// ToListPageNoticeLog 通知消息列表
// @post noticeLogList
// @filter application.Jwt
func ToListPageNoticeLog(appId string, pageSize, pageIndex int, monitorRepository monitor.Repository) collections.PageList[monitor.NoticeLogEO] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	return monitorRepository.ToListPageNoticeLog(appId, pageSize, pageIndex)
}

// DeleteNoticeLog 删除通知消息日志
// @post delNoticeLog
// @filter application.Jwt
func DeleteNoticeLog(startTime, endTime string, monitorRepository monitor.Repository) {
	stime := parse.ToTime(startTime)
	etime := parse.ToTime(endTime)
	err := monitorRepository.DeleteNoticeLog(stime, etime)
	exception.ThrowWebExceptionError(403, err)
}
