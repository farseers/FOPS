package request

import (
	"fmt"
	"fops/domain/_/eumBackupStoreType"
	"fops/domain/backupData"
	"strings"

	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/utils/cloud/aliyun"
	"github.com/farseer-go/webapi/check"
)

type AddRequest struct {
	GetDatabaseListRequest
	Database    []string                // 数据库
	Cron        string                  // 备份间隔
	StoreType   eumBackupStoreType.Enum // 备份存储类型
	StoreConfig string                  // 备份存储配置
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
		var ossStoreConfig aliyun.OSSConfig
		err := snc.Unmarshal([]byte(receiver.StoreConfig), &ossStoreConfig)
		check.IsTrue(err != nil, 403, fmt.Sprintf("OSS配置解析失败：%v", err))

		check.IsTrue(len(ossStoreConfig.AccessKeyID) == 0, 403, "AccessKeyID不能为空")
		check.IsTrue(len(ossStoreConfig.AccessKeySecret) == 0, 403, "AccessKeySecret不能为空")
		check.IsTrue(len(ossStoreConfig.BucketName) == 0, 403, "BucketName不能为空")

	case eumBackupStoreType.LocalDirectory:
		var fileStoreConfig backupData.FileStoreConfig
		err := snc.Unmarshal([]byte(receiver.StoreConfig), &fileStoreConfig)
		check.IsTrue(err != nil, 403, fmt.Sprintf("目录配置解析失败：%v", err))

		check.IsTrue(len(fileStoreConfig.Directory) == 0, 403, "目录不能为空")
	}
}
