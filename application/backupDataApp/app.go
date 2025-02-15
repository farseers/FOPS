// @area /backupData/
package backupDataApp

import (
	"fmt"
	"fops/application/backupDataApp/request"
	"fops/domain/_/eumBackupDataType"
	"fops/domain/_/eumBackupStoreType"
	"fops/domain/backupData"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/webapi/check"
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
	do.NextBackupAt = dateTime.Now() // 先改成立即备份
	// 生成ID
	do.GenerateId()
	count := backupDataRepository.GetCountById(do.Id)
	if count > 0 {
		do.Id += "_" + parse.ToString(count)
	}

	if do.StoreType == eumBackupStoreType.OSS {
		ossConfig, err := do.GetOSSConfig()
		if err != nil {
			flog.Warning(err.Error())
		}
		client, _, err := ossConfig.GetOssClient()
		check.IsTrue(client == nil || err != nil, 403, "OSS尝试连接失败，请确认鉴权是否正确:"+err.Error())
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
	lst := backupDataRepository.ToList()
	lst.Foreach(func(item *backupData.DomainObject) {
		item.Password = ""
	})
	return lst
}

// Info 备份计划查询
// @post info
// @filter application.Jwt
func Info(id string, backupDataRepository backupData.Repository) backupData.DomainObject {
	do := backupDataRepository.ToEntity(id)
	check.IsTrue(do.IsNil(), 403, "备份计划不存在")

	do.Password = ""
	return do
}

// Backup 立即备份
// @post backup
// @filter application.Jwt
func Backup(id string, backupDataRepository backupData.Repository) {
	do := backupDataRepository.ToEntity(id)
	check.IsTrue(do.IsNil(), 403, "备份计划不存在")

	if err := do.Backup(); err != nil {
		exception.ThrowRefuseExceptionError(err)
	}
	backupDataRepository.UpdateAt(do.Id, do)
}

// Delete 删除备份计划
// @post delete
// @filter application.Jwt
func Delete(id string, backupDataRepository backupData.Repository) {
	do := backupDataRepository.ToEntity(id)
	check.IsTrue(do.IsNil(), 403, "备份计划不存在")

	// 遍历数据库，删除备份文件
	for _, database := range do.Database {
		lstHistoryData, err := do.GetHistoryData(database)
		for lstHistoryData.Count() > 0 && err == nil {
			lstHistoryData.Foreach(func(item *backupData.BackupHistoryData) {
				do.DeleteBackupFile(item.FileName)
			})
			lstHistoryData, err = do.GetHistoryData(database)
		}
	}

	// 删除备份计划
	backupDataRepository.Delete(id)
}

// Clear 清空备份计划中的所有备份文件
// @post clear
// @filter application.Jwt
func Clear(id string, backupDataRepository backupData.Repository) {
	do := backupDataRepository.ToEntity(id)
	check.IsTrue(do.IsNil(), 403, "备份计划不存在")

	// 遍历数据库，删除备份文件
	for _, database := range do.Database {
		lstHistoryData, err := do.GetHistoryData(database)
		for lstHistoryData.Count() > 0 && err == nil {
			lstHistoryData.Foreach(func(item *backupData.BackupHistoryData) {
				do.DeleteBackupFile(item.FileName)
			})
			lstHistoryData, err = do.GetHistoryData(database)
		}
	}
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
func BackupList(backupId string, database string, backupDataRepository backupData.Repository) collections.List[backupData.BackupHistoryData] {
	do := backupDataRepository.ToEntity(backupId)
	check.IsTrue(do.IsNil(), 403, "备份计划不存在")

	lst, err := do.GetHistoryData(database)
	exception.ThrowRefuseExceptionError(err)

	return lst
}

// DeleteHistory 删除备份文件
// @post deleteHistory
// @filter application.Jwt
func DeleteBackupFile(backupId string, fileName string, backupDataRepository backupData.Repository) {
	do := backupDataRepository.ToEntity(backupId)
	check.IsTrue(do.IsNil(), 403, "备份计划不存在")

	err := do.DeleteBackupFile(fileName)
	exception.ThrowRefuseExceptionError(err)
}

// RecoverBackupFile 恢复备份文件
// @post recoverBackupFile
// @filter application.Jwt
func RecoverBackupFile(backupId string, database string, fileName string, backupDataRepository backupData.Repository) {
	do := backupDataRepository.ToEntity(backupId)
	check.IsTrue(do.IsNil(), 403, "备份计划不存在")

	err := do.RecoverBackupFile(database, fileName)
	exception.ThrowRefuseExceptionError(err)
}
