package repository

import (
	"fops/domain/backupData"
	"fops/infrastructure/repository/context"
	"fops/infrastructure/repository/model"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/mapper"
)

type backupDataRepository struct {
	// IRepository 通用的仓储接口
	data.IRepository[backupData.DomainObject]
}

// 查看同名的ID的数量
func (receiver *backupDataRepository) GetCountById(id string) int64 {
	return context.MysqlContext.BackupData.Where("id = ?", id).Count()
}

// 修改备份计划
func (receiver *backupDataRepository) Update(id any, do backupData.DomainObject) (int64, error) {
	po := mapper.Single[model.BackupDataPO](do)
	return context.MysqlContext.BackupData.Where("id = ?", id).Omit("id", "backup_data_type", "store_type", "last_backup_at").Update(po)
}

// 更新备份计划的时间
func (receiver *backupDataRepository) UpdateAt(id any, do backupData.DomainObject) (int64, error) {
	po := mapper.Single[model.BackupDataPO](do)
	return context.MysqlContext.BackupData.Where("id = ?", id).Select("last_backup_at", "next_backup_at").Update(po)
}

// 获取即将备份的数据
func (receiver *backupDataRepository) ToNextBackupData() backupData.DomainObject {
	po := context.MysqlContext.BackupData.Asc("next_backup_at").ToEntity()
	return mapper.Single[backupData.DomainObject](po)
}

// 添加备份文件列表
func (receiver *backupDataRepository) AddHistory(lst collections.List[backupData.BackupHistoryData]) {
	lstPO := mapper.ToList[model.BackupHistoryDataPO](lst)
	context.MysqlContext.BackupHistoryData.InsertList(lstPO, 2000)
}

// 删除备份文件
func (receiver *backupDataRepository) DeleteHistory(backupId string, fileName string) (int64, error) {
	return context.MysqlContext.BackupHistoryData.Where("backup_id = ? and file_name = ?", backupId, fileName).Delete()
}

// 获取备份文件列表
func (receiver *backupDataRepository) ToBackupList(backupId string) collections.List[backupData.BackupHistoryData] {
	lstPO := context.MysqlContext.BackupHistoryData.Where("backup_id = ?", backupId).Desc("create_at").Limit(144).ToList()
	return mapper.ToList[backupData.BackupHistoryData](lstPO)
}
