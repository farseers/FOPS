package backupData

import (
	"context"
	"fmt"
	"fops/domain/_/eumBackupDataType"
	"fops/domain/_/eumBackupStoreType"
	"fops/domain/apps"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snc"
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
	LastBackupAt   dateTime.DateTime       // 上次备份时间
	NextBackupAt   dateTime.DateTime       // 下次备份时间
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
func (receiver *DomainObject) Backup() {
	var lstBackupHistoryData collections.List[BackupHistoryData]

	// 确定本地存储目录
	backupRoot := receiver.getBackupRoot()

	// 备份数据
	switch receiver.BackupDataType {
	case eumBackupDataType.Mysql:
		lstBackupHistoryData = receiver.backupMySQL(backupRoot)
	}

	if lstBackupHistoryData.Count() == 0 {
		return
	}

	// 上传备份文件
	if receiver.StoreType == eumBackupStoreType.OSS {
		receiver.uploadOSS(lstBackupHistoryData)
	}
}

// 备份MySQL
func (receiver *DomainObject) backupMySQL(backupRoot string) collections.List[BackupHistoryData] {
	// 安装 mysqldump
	if !isMysqldumpInstalled() {
		installMysqldump()
	}

	lstBackupHistoryData := collections.NewList[BackupHistoryData]()
	// 备份数据库
	for _, database := range receiver.Database {
		filePath := receiver.Id + "/" + database + "/"
		// 创建备份目录
		if !file.IsExists(backupRoot + filePath) {
			file.CreateDir766(backupRoot + filePath)
		}

		fileName := filePath + database + "_" + time.Now().Format("2006_01_02_15_04") + ".sql.gz"
		mysqldumpCmd := fmt.Sprintf("mysqldump -h %s -P %d -u%s -p%s %s | gzip > %s", receiver.Host, receiver.Port, receiver.Username, receiver.Password, database, backupRoot+fileName)
		code, result := exec.RunShellCommand(mysqldumpCmd, nil, "", false)
		// 备份失败时删除备份文件
		if code != 0 {
			file.Delete(backupRoot + fileName)
			flog.Warningf("备份%s数据库失败：%s", database, collections.NewList(result...).ToString(","))
			continue
		}
		fileInfo, err := os.Stat(backupRoot + fileName)
		if err != nil {
			flog.Warningf("获取备份文件信息:%s,失败： %s", backupRoot+fileName, err.Error())
			continue
		}
		lstBackupHistoryData.Add(BackupHistoryData{
			BackupId:  receiver.Id,
			Database:  database,
			FileName:  fileName,
			StoreType: receiver.StoreType,
			CreateAt:  dateTime.Now(),
			Size:      fileInfo.Size() / 1024,
		})
	}
	return lstBackupHistoryData
}

// 获取OSS客户端
func (receiver *DomainObject) getOssClient() (*oss.Client, string) {
	var ossStoreConfig OSSStoreConfig
	err := snc.Unmarshal([]byte(receiver.StoreConfig), &ossStoreConfig)
	if err != nil {
		flog.Warningf("OSS的配置解析失败：%v", err)
		return nil, ""
	}

	// 从环境变量中获取访问凭证。运行本代码示例之前，请确保已设置环境变量OSS_ACCESS_KEY_ID和OSS_ACCESS_KEY_SECRET。
	cfg := oss.LoadDefaultConfig().WithCredentialsProvider(credentials.NewStaticCredentialsProvider(ossStoreConfig.AccessKeyID, ossStoreConfig.AccessKeySecret))
	if ossStoreConfig.Region != "" {
		cfg.WithRegion(ossStoreConfig.Region)
	}
	if ossStoreConfig.Endpoint != "" {
		cfg.WithEndpoint(ossStoreConfig.Endpoint)
	}

	return oss.NewClient(cfg), ossStoreConfig.BucketName
}

// 上传备份文件到OSS
func (receiver *DomainObject) uploadOSS(lstBackupHistoryData collections.List[BackupHistoryData]) {
	client, bucketName := receiver.getOssClient()
	if client == nil {
		return
	}
	backupRoot := receiver.getBackupRoot()
	// 批量上传
	for _, item := range lstBackupHistoryData.ToArray() {
		f, err := os.Open(backupRoot + item.FileName)
		if err != nil {
			flog.Warningf("打开上传文件：%s 时，发生错误：%v", item.FileName, err)
			//lstBackupHistoryData.RemoveAt(index)
			continue
		}
		defer f.Close()

		result, err := client.PutObject(context.TODO(), &oss.PutObjectRequest{
			Bucket: oss.Ptr(bucketName),
			Key:    oss.Ptr(item.FileName),
			Body:   f,
		})

		if err != nil {
			flog.Warningf("上传文件：%s 时，发生错误：%v", item.FileName, err)
			//lstBackupHistoryData.RemoveAt(index)
		}

		// 上传成功后，删除本地文件
		file.Delete(backupRoot + item.FileName)
		flog.Infof("数据库：%s，OSS上传文件：%s 成功, ETag :%v\n", item.Database, item.FileName, result.ETag)
	}
}

// 删除备份文件
func (receiver *DomainObject) DeleteBackupFile(fileName string) {
	// 删除本地文件
	file.Delete(receiver.getBackupRoot() + fileName)
	// 删除OSS文件
	if receiver.StoreType == eumBackupStoreType.OSS {
		client, bucketName := receiver.getOssClient()
		if client == nil {
			return
		}

		// 执行删除对象的操作并处理结果
		result, err := client.DeleteObject(context.TODO(), &oss.DeleteObjectRequest{
			Bucket: oss.Ptr(bucketName), // 存储空间名称
			Key:    oss.Ptr(fileName),   // 对象名称
		})
		if err != nil {
			flog.Warningf("OSS删除文件：%s 时，发生错误：%v", fileName, err)
		} else {
			// 打印删除对象的结果
			flog.Infof("OSS删除文件：%s 成功, ETag :%v\n", fileName, result)
		}
	}
}

// 恢复备份文件
func (receiver *DomainObject) RecoverBackupFile(fileName string) {

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
func (receiver *DomainObject) GetHistoryData(database string) collections.List[BackupHistoryData] {
	filePath := receiver.Id + "/" + database + "/"

	lstBackupHistoryData := collections.NewList[BackupHistoryData]()
	// 通过oss获取
	switch receiver.StoreType {
	case eumBackupStoreType.OSS:
		client, bucketName := receiver.getOssClient()
		if client == nil {
			return lstBackupHistoryData
		}

		// 执行列举所有文件的操作
		lsRes, err := client.ListObjectsV2(context.TODO(), &oss.ListObjectsV2Request{
			Bucket:            oss.Ptr(bucketName),
			ContinuationToken: oss.Ptr(""),
			Prefix:            oss.Ptr(filePath), // 列举指定目录下的所有对象
			MaxKeys:           144,
		})

		if err != nil {
			flog.Warningf("OSS读取文件列表：%s 时，发生错误：%v", database, err)
		}

		// 打印列举结果
		for _, object := range lsRes.Contents {
			lstBackupHistoryData.Add(BackupHistoryData{
				BackupId:  receiver.Id,
				Database:  database,
				FileName:  *object.Key,
				StoreType: receiver.StoreType,
				CreateAt:  dateTime.New(*object.LastModified),
				Size:      object.Size / 1024,
			})
			log.Printf("Object Key: %s, Type: %s, Size: %d, ETag: %s, LastModified: %s, StorageClass: %s\n",
				*object.Key, *object.Type, object.Size, *object.ETag, object.LastModified.Format(time.RFC3339), *object.StorageClass)
		}
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
	Endpoint        string // 填写Bucket对应的Endpoint，以华东1（杭州）为例，填写为https://oss-cn-hangzhou.aliyuncs.com。其它Region请按实际情况填写。
	Region          string // 填写Bucket所在地域，以华东1（杭州）为例，填写为cn-hangzhou。其它Region请按实际情况填写。
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
	Database  string                  // 数据库
	StoreType eumBackupStoreType.Enum // 备份存储类型
	CreateAt  dateTime.DateTime       // 备份时间
	Size      int64                   // 备份文件大小（KB）
}
