package backupData

import (
	"fops/domain/_/eumBackupDataType"
	"fops/domain/_/eumBackupStoreType"
	"time"

	"github.com/farseer-go/fs/parse"
)

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
	Id        string                  // 备份计划的ID
	FileName  string                  // 文件名
	StoreType eumBackupStoreType.Enum // 备份存储类型
	CreateAt  int64                   // 备份时间
	Size      int64                   // 备份文件大小（KB）
}
