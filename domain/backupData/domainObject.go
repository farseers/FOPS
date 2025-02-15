package backupData

import (
	"fmt"
	"fops/domain/_/eumBackupDataType"
	"fops/domain/_/eumBackupStoreType"
	"fops/domain/apps"
	"os"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/utils/cloud/aliyun"
	"github.com/farseer-go/utils/db"
	"github.com/farseer-go/utils/file"
	"github.com/robfig/cron/v3"
)

var StandardParser = cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

// DomainObject 备份计划
type DomainObject struct {
	Id             string                  // ID
	BackupDataType eumBackupDataType.Enum  // 备份数据类型
	Host           string                  // 主机
	Port           int                     // 端口
	Username       string                  // 用户名
	Password       string                  // 密码
	Database       []string                // 数据库
	LastBackupAt   dateTime.DateTime       // 上次备份时间
	NextBackupAt   dateTime.DateTime       // 下次备份时间
	Cron           string                  // 备份间隔
	StoreType      eumBackupStoreType.Enum // 备份存储类型
	StoreConfig    string                  // 备份存储配置
}

// 生成ID
func (receiver *DomainObject) GenerateId() {
	receiver.Id = receiver.Host + "_" + parse.ToString(receiver.Port) + "_" + receiver.Username
}

// 获取OSS配置
func (receiver *DomainObject) GetOSSConfig() (*aliyun.OSSConfig, error) {
	var ossConfig aliyun.OSSConfig
	err := snc.Unmarshal([]byte(receiver.StoreConfig), &ossConfig)
	if err != nil {
		return nil, fmt.Errorf("OSS的配置解析失败：%v", err)
	}
	return &ossConfig, nil
}

func (receiver *DomainObject) IsNil() bool {
	return receiver.Id == ""
}

// 备份
func (receiver *DomainObject) Backup() error {
	// 更新时间，放在前面处理是节省后面多个return时的逻辑处理，这里更简单
	cornSchedule, err := StandardParser.Parse(receiver.Cron)
	if err != nil {
		return fmt.Errorf("同步备份计划时，do.Cron的值不正确导致错误: %s %v", receiver.Cron, err)
	}
	// 更新时间字段，并生成下一次执行时间。
	receiver.LastBackupAt = dateTime.Now()
	receiver.NextBackupAt = dateTime.New(cornSchedule.Next(time.Now()))

	var lstBackupHistoryData collections.List[BackupHistoryData]
	// 确定本地存储目录
	backupRoot := receiver.getBackupRoot()

	// 备份数据
	switch receiver.BackupDataType {
	case eumBackupDataType.Mysql:
		lstBackupHistoryData = receiver.backupMySQL(backupRoot)
	}

	if lstBackupHistoryData.Count() == 0 {
		return nil
	}

	// 上传备份文件
	if receiver.StoreType == eumBackupStoreType.OSS {
		ossConfig, err := receiver.GetOSSConfig()
		if err != nil {
			return err
		}

		var fileNames []string
		lstBackupHistoryData.Select(&fileNames, func(item BackupHistoryData) any {
			return item.FileName
		})
		return ossConfig.UploadOSS(backupRoot, fileNames)
	}
	return nil
}

// 备份MySQL
func (receiver *DomainObject) backupMySQL(backupRoot string) collections.List[BackupHistoryData] {
	lstBackupHistoryData := collections.NewList[BackupHistoryData]()
	// 备份数据库
	for _, database := range receiver.Database {
		filePath := receiver.Id + "/" + database + "/"
		// 创建备份目录
		if !file.IsExists(backupRoot + filePath) {
			file.CreateDir766(backupRoot + filePath)
		}

		fileName := filePath + database + "_" + time.Now().Format("2006_01_02_15_04") + ".sql.gz"
		fileSize, err := db.BackupMysql(receiver.Host, receiver.Port, receiver.Username, receiver.Password, database, backupRoot+fileName)
		if err != nil {
			flog.Warning(err.Error())
			continue
		}

		lstBackupHistoryData.Add(BackupHistoryData{
			BackupId:  receiver.Id,
			Database:  database,
			FileName:  fileName,
			StoreType: receiver.StoreType,
			CreateAt:  dateTime.Now(),
			Size:      fileSize,
		})
	}
	return lstBackupHistoryData
}

// 删除备份文件
func (receiver *DomainObject) DeleteBackupFile(fileName string) error {
	// 删除本地文件
	file.Delete(receiver.getBackupRoot() + fileName)
	// 删除OSS文件
	if receiver.StoreType == eumBackupStoreType.OSS {
		ossConfig, err := receiver.GetOSSConfig()
		if err != nil {
			return err
		}
		err = ossConfig.DeleteFile(fileName)
		if err != nil {
			return err
		}
	}
	return nil
}

// 恢复备份文件
func (receiver *DomainObject) RecoverBackupFile(database string, fileName string) error {
	// 确定本地存储目录
	backupRoot := receiver.getBackupRoot()

	switch receiver.StoreType {
	// 通过oss获取，需先下载文件到本地目录
	case eumBackupStoreType.OSS:
		ossConfig, err := receiver.GetOSSConfig()
		if err != nil {
			flog.Warning(err.Error())
		}
		err = ossConfig.DownloadFile(backupRoot, fileName)
		if err != nil {
			return err
		}
	}

	// 恢复数据库
	err := db.RecoverMysql(receiver.Host, receiver.Port, receiver.Username, receiver.Password, database, backupRoot+fileName)
	// 如果是oss下载的，则删除原文件
	if receiver.StoreType == eumBackupStoreType.OSS {
		file.Delete(backupRoot + fileName)
	}
	return err
}

// 备份的根目录
func (receiver *DomainObject) getBackupRoot() string {
	backupRoot := apps.BackupRoot
	if receiver.StoreType == eumBackupStoreType.LocalDirectory {
		var fileStoreConfig FileStoreConfig
		err := snc.Unmarshal([]byte(receiver.StoreConfig), &fileStoreConfig)
		if err != nil {
			flog.Warningf("备份%s时，目录配置解析失败：%v", receiver.Id, err)
		} else {
			backupRoot = fileStoreConfig.Directory
		}
	}
	return backupRoot
}

// 获取备份历史数据
func (receiver *DomainObject) GetHistoryData(database string) (collections.List[BackupHistoryData], error) {
	filePath := receiver.Id + "/" + database + "/"

	lstBackupHistoryData := collections.NewList[BackupHistoryData]()
	// 通过oss获取
	switch receiver.StoreType {
	case eumBackupStoreType.OSS:
		ossConfig, err := receiver.GetOSSConfig()
		if err != nil {
			return lstBackupHistoryData, err
		}
		fileObjects, err := ossConfig.GetFileList(filePath)
		if err != nil {
			return lstBackupHistoryData, err
		}
		fileObjects.Foreach(func(item *aliyun.FileObject) {
			lstBackupHistoryData.Add(BackupHistoryData{
				BackupId:  receiver.Id,
				Database:  database,
				FileName:  item.FileName,
				StoreType: receiver.StoreType,
				CreateAt:  item.CreateAt,
				Size:      item.Size,
			})
		})
	case eumBackupStoreType.LocalDirectory:
		filePath = receiver.getBackupRoot() + filePath
		lst := file.GetFiles(filePath, "*", true)
		for _, file := range lst {
			fileInfo, _ := os.Stat(file)
			if fileInfo != nil {
				lstBackupHistoryData.Add(BackupHistoryData{
					BackupId:  receiver.Id,
					Database:  database,
					FileName:  file[len(receiver.getBackupRoot()):],
					StoreType: receiver.StoreType,
					CreateAt:  dateTime.New(fileInfo.ModTime()),
					Size:      fileInfo.Size() / 1024 / 1024,
				})
			}
		}
	}
	return lstBackupHistoryData, nil
}

// 本地目录存储配置
type FileStoreConfig struct {
	Directory string // 目录
}

// 备份历史数据
type BackupHistoryData struct {
	BackupId  string                  // 备份计划的ID
	FileName  string                  // 文件名
	Database  string                  // 数据库
	StoreType eumBackupStoreType.Enum // 备份存储类型
	CreateAt  dateTime.DateTime       // 备份时间
	Size      int64                   // 备份文件大小（KB）
}
