package backupData

import (
	"fmt"
	"fops/domain/_/eumBackupDataType"
	"fops/domain/_/eumBackupStoreType"
	"fops/domain/apps"
	"os"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/utils/exec"
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
	LastBackupAt   time.Time               // 上次备份时间
	NextBackupAt   time.Time               // 下次备份时间
	Cron           string                  // 备份间隔
	StoreType      eumBackupStoreType.Enum // 备份存储类型
	StoreConfig    string                  // 备份存储配置
}

func (receiver *DomainObject) GenerateId() {
	receiver.Id = receiver.Host + "_" + parse.ToString(receiver.Port) + "_" + receiver.Username
}

func (receiver *DomainObject) IsNil() bool {
	return receiver.Id == ""
}

// 备份
func (receiver *DomainObject) Backup() collections.List[BackupHistoryData] {
	var lstBackupHistoryData collections.List[BackupHistoryData]

	// 备份数据
	switch receiver.BackupDataType {
	case eumBackupDataType.Mysql:
		lstBackupHistoryData = receiver.backupMySQL()
	}

	if lstBackupHistoryData.Count() == 0 {
		return lstBackupHistoryData
	}

	// 上传备份文件
	switch receiver.StoreType {
	case eumBackupStoreType.OSS:
		//ossStoreConfig := mapper.Single[OSSStoreConfig](receiver.StoreConfig)

	case eumBackupStoreType.LocalDirectory:
		//fileStoreConfig := mapper.Single[FileStoreConfig](receiver.StoreConfig)
	}
	return lstBackupHistoryData
}

// 备份MySQL
func (receiver *DomainObject) backupMySQL() collections.List[BackupHistoryData] {
	// 安装 mysqldump
	if !isMysqldumpInstalled() {
		installMysqldump()
	}

	lstBackupHistoryData := collections.NewList[BackupHistoryData]()
	// 备份数据库
	for _, database := range receiver.Database {
		flog.Info("3.1 ")
		filePath := apps.BackupRoot + receiver.Host + "_" + database + "_$(date +%Y%m%d_%H_%M).sql.gz"
		mysqldumpCmd := fmt.Sprintf("mysqldump -h %s -P %d -u%s -p%s %s | gzip > %s", receiver.Host, receiver.Port, receiver.Username, receiver.Password, database, filePath)
		code, result := exec.RunShellCommand(mysqldumpCmd, nil, "", false)
		flog.Info("3.2 ")
		// 备份失败时删除备份文件
		if code != 0 {
			flog.Info("3.3 ")
			file.Delete(filePath)
			flog.Warningf("备份%s数据库失败：%s", database, collections.NewList(result...).ToString(","))
			continue
		}
		flog.Info("3.4 ")
		fileInfo, _ := os.Stat(filePath)
		lstBackupHistoryData.Add(BackupHistoryData{
			BackupId:  receiver.Id,
			FileName:  filePath,
			StoreType: receiver.StoreType,
			CreateAt:  time.Now(),
			Size:      fileInfo.Size() / 1024,
		})
	}
	return lstBackupHistoryData
}

// 检查 mysqldump 是否已安装
func isMysqldumpInstalled() bool {
	code, result := exec.RunShellCommand("mysqldump --version", nil, "", false)
	if code != 0 || len(result) == 0 {
		return false
	}
	// 检查输出中是否包含 "mysqldump" 关键字
	return strings.Contains(result[0], "mysqldump")
}

// 安装 mysqldump
func installMysqldump() {
	exec.RunShellCommand("apk add --no-cache mariadb-client", nil, "", false)
}

// 阿里云OSS存储配置
type OSSStoreConfig struct {
	AccessKeyID     string // AccessKeyID
	AccessKeySecret string // AccessKeySecret
	Endpoint        string // Endpoint
	BucketName      string // BucketName
}

// 本地目录存储配置
type FileStoreConfig struct {
	Directory string // 目录
}

// 备份历史数据
type BackupHistoryData struct {
	BackupId  string                  // 备份计划的ID
	FileName  string                  // 文件名
	StoreType eumBackupStoreType.Enum // 备份存储类型
	CreateAt  time.Time               // 备份时间
	Size      int64                   // 备份文件大小（KB）
}
