// @area /backupData/
package backupDataApp

import (
	"fmt"
	"fops/application/backupDataApp/request"
	"fops/domain/_/eumBackupDataType"
	"fops/domain/backupData"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/dateTime"
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
	do.LastBackupAt = dateTime.Now()
	do.NextBackupAt = dateTime.New(cornSchedule.Next(time.Now()))

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
	do.NextBackupAt = dateTime.New(cornSchedule.Next(time.Now()))

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

// Delete 删除备份计划
// @post delete
// @filter application.Jwt
func Delete(id string, backupDataRepository backupData.Repository) {
	do := backupDataRepository.ToEntity(id)
	lstHistoryData := backupDataRepository.ToBackupList(id)
	lstHistoryData.Foreach(func(item *backupData.BackupHistoryData) {
		do.DeleteBackupFile(item.FileName)
		backupDataRepository.DeleteHistory(id, item.FileName)
	})
}

// GetDatabaseList 获取数据库列表
// @post getDatabaseList
// @filter application.Jwt
func GetDatabaseList(req request.GetDatabaseListRequest) []string {
	var dbConnectionString string
	switch req.BackupDataType {
	case eumBackupDataType.Mysql:
		dbConnectionString = fmt.Sprintf("DataType=mysql,ConnectionString=%s:%s@tcp(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", req.Username, req.Password, req.Host, req.Port)
	case eumBackupDataType.Clickhouse:
		dbConnectionString = fmt.Sprintf("DataType=clickhouse,ConnectionString=tcp://%s:%d?username=%s&password=%s&read_timeout=10&write_timeout=20", req.Host, req.Port, req.Username, req.Password)
	}
	databases, err := data.NewInternalContext(dbConnectionString).GetDatabaseList()
	exception.ThrowRefuseExceptionError(err)
	return databases
}

// BackupList 备份文件列表
// @post backupList
// @filter application.Jwt
func BackupList(backupId string, backupDataRepository backupData.Repository) collections.List[backupData.BackupHistoryData] {
	return backupDataRepository.ToBackupList(backupId)
}

// DeleteHistory 删除备份文件
// @post deleteHistory
// @filter application.Jwt
func DeleteBackupFile(backupId string, fileName string, backupDataRepository backupData.Repository) {
	do := backupDataRepository.ToEntity(backupId)
	do.DeleteBackupFile(fileName)
	backupDataRepository.DeleteHistory(backupId, fileName)
}

// RecoverBackupFile 恢复备份文件
// @post recoverBackupFile
// @filter application.Jwt
func RecoverBackupFile(backupId string, fileName string, backupDataRepository backupData.Repository) {
	do := backupDataRepository.ToEntity(backupId)
	do.RecoverBackupFile(fileName)
}
