// @area /backupData/
package backupDataApp

import (
	"fops/application/backupDataApp/request"
	"fops/domain/backupData"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/mapper"
)

// Add 添加备份计划
// @post add
// @filter application.Jwt
func Add(req request.AddRequest, backupDataRepository backupData.Repository) {
	do := mapper.Single[backupData.DomainObject](req)

	// 验证cron
	cornSchedule, err := backupData.StandardParser.Parse(do.Cron)
	if err != nil {
		exception.ThrowWebExceptionf(403, "Cron格式错误:%s", do.Cron)
	}
	do.NextBackupAt = cornSchedule.Next(time.Now())

	// 生成ID
	do.GenerateId()
	count := backupDataRepository.GetCountById(do.Id)
	if count > 0 {
		do.Id += "_" + parse.ToString(count)
	}

	// 添加
	err = backupDataRepository.Add(do)
	exception.ThrowWebExceptionError(403, err)
}

// Update 更新备份计划
// @post update
// @filter application.Jwt
func Update(req request.UpdateRequest, backupDataRepository backupData.Repository) {
	do := mapper.Single[backupData.DomainObject](req)

	// 验证cron
	cornSchedule, err := backupData.StandardParser.Parse(do.Cron)
	if err != nil {
		exception.ThrowWebExceptionf(403, "Cron格式错误:%s", do.Cron)
	}
	do.NextBackupAt = cornSchedule.Next(time.Now())

	// 修改
	_, err = backupDataRepository.Update(req.Id, do)
	exception.ThrowWebExceptionError(403, err)
}

// List 备份计划列表
// @post list
// @filter application.Jwt
func List(backupDataRepository backupData.Repository) collections.List[backupData.DomainObject] {
	return backupDataRepository.ToList()
}

// Info 备份计划查询
// @post info
// @filter application.Jwt
func Info(id string, backupDataRepository backupData.Repository) backupData.DomainObject {
	return backupDataRepository.ToEntity(id)
}
