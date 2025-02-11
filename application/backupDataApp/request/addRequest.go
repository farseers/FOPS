package request

import (
	"fops/domain/_/eumBackupDataType"
	"fops/domain/_/eumBackupStoreType"
	"fops/domain/backupData"
	"strings"

	"github.com/farseer-go/mapper"
	"github.com/farseer-go/webapi/check"
)

type AddRequest struct {
	BackupDataType eumBackupDataType.Enum  // 备份数据类型
	Host           string                  // 主机
	Port           int                     // 端口
	Username       string                  // 用户名
	Password       string                  // 密码
	Database       []string                // 数据库
	Cron           string                  // 备份间隔
	StoreType      eumBackupStoreType.Enum // 备份存储类型
	StoreConfig    string                  // 备份存储配置
}

func (receiver *AddRequest) Check() {
	// 去除空格
	receiver.Host = strings.TrimSpace(receiver.Host)
	receiver.Username = strings.TrimSpace(receiver.Username)
	receiver.Password = strings.TrimSpace(receiver.Password)
	receiver.Cron = strings.TrimSpace(receiver.Cron)
	receiver.StoreConfig = strings.TrimSpace(receiver.StoreConfig)

	// 主机
	check.IsTrue(len(receiver.Host) == 0, 403, "主机不能为空")
	// 端口
	check.IsTrue(receiver.Port < 1, 403, "端口不能为空")
	// 用户名
	check.IsTrue(len(receiver.Username) == 0, 403, "用户名不能为空")
	// 数据库
	check.IsTrue(len(receiver.Database) == 0, 403, "数据库不能为空")
	// 备份间隔
	check.IsTrue(len(receiver.Cron) == 0, 403, "备份间隔不能为空")
	// 备份存储配置
	check.IsTrue(len(receiver.StoreConfig) == 0, 403, "备份存储配置不能为空")

	switch receiver.StoreType {
	case eumBackupStoreType.OSS:
		oSSStoreConfig := mapper.Single[backupData.OSSStoreConfig](receiver.StoreConfig)
		check.IsTrue(len(oSSStoreConfig.AccessKeyID) == 0, 403, "AccessKeyID不能为空")
		check.IsTrue(len(oSSStoreConfig.AccessKeySecret) == 0, 403, "AccessKeySecret不能为空")
		check.IsTrue(len(oSSStoreConfig.BucketName) == 0, 403, "BucketName不能为空")

	case eumBackupStoreType.LocalDirectory:
		fileStoreConfig := mapper.Single[backupData.FileStoreConfig](receiver.StoreConfig)
		check.IsTrue(len(fileStoreConfig.Directory) == 0, 403, "目录不能为空")
	}
}
