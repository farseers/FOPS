package backupData

import (
	"bufio"
	"fmt"
	"fops/domain/_/eumBackupDataType"
	"fops/domain/_/eumBackupStoreType"
	"fops/domain/apps"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/data"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/flog"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/fs/snc"
	"github.com/farseer-go/fs/stopwatch"
	"github.com/farseer-go/utils/cloud/aliyun"
	"github.com/farseer-go/utils/db"
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

	// 确定本地存储目录
	backupRoot := receiver.getBackupRoot()

	lstErrorContent := collections.NewList[string]()
	for _, database := range receiver.Database {
		var lstBackupHistoryData collections.List[BackupHistoryData]
		// 备份数据
		switch receiver.BackupDataType {
		case eumBackupDataType.Mysql:
			lstBackupHistoryData, err = receiver.backupMySQL(backupRoot, database)
		case eumBackupDataType.Clickhouse:
			lstBackupHistoryData, err = receiver.backupClickhouse(backupRoot, database)
		}

		if err != nil {
			lstErrorContent.Add(fmt.Sprintf("在备份%s时，发生错误：%v", database, err))
			continue
		}
		if lstBackupHistoryData.Count() == 0 {
			continue
		}

		// 上传备份文件
		if receiver.StoreType == eumBackupStoreType.OSS {
			ossConfig, err := receiver.GetOSSConfig()
			if err != nil {
				lstErrorContent.Add(fmt.Sprintf("在上传备份文件%s时，发生错误：%v", database, err))
				continue
			}

			var fileNames []string
			lstBackupHistoryData.Select(&fileNames, func(item BackupHistoryData) any {
				return item.FileName
			})
			err = ossConfig.UploadOSS(backupRoot, fileNames)
			if err != nil {
				lstErrorContent.Add(fmt.Sprintf("在上传备份文件%s时，发生错误：%v", database, err))
				continue
			}
		}
	}

	if lstErrorContent.Count() > 0 {
		return fmt.Errorf("%s", lstErrorContent.ToString("\n"))
	}
	return nil
}

// 备份MySQL
func (receiver *DomainObject) backupMySQL(backupRoot, database string) (collections.List[BackupHistoryData], error) {
	lstBackupHistoryData := collections.NewList[BackupHistoryData]()
	filePath := receiver.Id + "/" + database + "/"
	// 创建备份目录
	if !file.IsExists(backupRoot + filePath) {
		file.CreateDir766(backupRoot + filePath)
	}

	fileName := filePath + database + "_" + time.Now().Format("2006_01_02_15_04") + ".sql.gz"
	fileSize, err := db.BackupMysql(receiver.Host, receiver.Port, receiver.Username, receiver.Password, database, backupRoot+fileName)
	if err != nil {
		return lstBackupHistoryData, err
	}

	lstBackupHistoryData.Add(BackupHistoryData{
		BackupId:  receiver.Id,
		Database:  database,
		FileName:  fileName,
		StoreType: receiver.StoreType,
		CreateAt:  dateTime.Now(),
		Size:      fileSize,
	})
	return lstBackupHistoryData, nil
}

// 备份Clickhouse
func (receiver *DomainObject) backupClickhouse(backupRoot, database string) (collections.List[BackupHistoryData], error) {
	lstBackupHistoryData := collections.NewList[BackupHistoryData]()

	// 备份数据库
	filePath := receiver.Id + "/" + database + "/"
	// 创建备份目录
	if !file.IsExists(backupRoot + filePath) {
		file.CreateDir766(backupRoot + filePath)
	}

	fileName := filePath + database + "_" + time.Now().Format("2006_01_02_15_04") + ".sql"
	path := filepath.Dir(backupRoot + fileName)
	file.CreateDir766(path)

	// 读取表列表
	dbConnectionString := data.CreateConnectionString(receiver.BackupDataType.ToString(), receiver.Host, receiver.Port, database, receiver.Username, receiver.Password)
	dbContext := data.NewInternalContext(dbConnectionString)
	tables, err := dbContext.GetTableList(database)
	if err != nil {
		return lstBackupHistoryData, err
	}

	// 删除文件（如果有）
	file.Delete(backupRoot + fileName)
	file.WriteString(backupRoot+fileName, "")

	for _, tableName := range tables {
		// 删除表
		file.AppendLine(backupRoot+fileName, fmt.Sprintf("DROP TABLE IF EXISTS %s.%s;", database, tableName))

		// 创建表
		var createTableSql string
		_, err = dbContext.ExecuteSqlToValue(&createTableSql, fmt.Sprintf("SHOW CREATE TABLE %s.%s", database, tableName))
		if err != nil {
			return lstBackupHistoryData, err
		}
		file.AppendLine(backupRoot+fileName, createTableSql+";")

		// 找到排序键
		var orderBy string
		re := regexp.MustCompile(`ORDER BY \((.*?)\)`)
		matches := re.FindStringSubmatch(createTableSql)
		if len(matches) > 0 {
			orderBy = matches[0]
		}

		// 得到总的数据量（用于分页计算）
		var totalCount float64
		_, err = dbContext.ExecuteSqlToValue(&totalCount, fmt.Sprintf("SELECT count() FROM %s.%s;", database, tableName))
		if err != nil {
			return lstBackupHistoryData, err
		}

		// 导出数据
		var realTotalCount float64
		pageSize := float64(10000)
		pageCount := math.Ceil(totalCount / pageSize)
		for pageIndex := float64(1); pageIndex <= pageCount; pageIndex++ {
			// 第二页起，休眠500ms
			if pageIndex > 1 {
				time.Sleep(500 * time.Millisecond)
			}

			sw := stopwatch.StartNew()

			offset := (pageIndex - 1) * pageSize

			var results []string
			query := fmt.Sprintf("SELECT toString(tuple(*)) FROM %s.%s %s LIMIT %d OFFSET %d FORMAT SQLInsert;", database, tableName, orderBy, int(pageSize), int(offset))
			_, err := dbContext.ExecuteSqlToResult(&results, query)
			if err != nil {
				return lstBackupHistoryData, err
			}

			realTotalCount += float64(len(results))

			insertSql := fmt.Sprintf("INSERT INTO %s.%s VALUES %s;", database, tableName, collections.NewList(results...).ToString(",\n"))
			file.AppendLine(backupRoot+fileName, insertSql)

			flog.Infof("导出%s.%s 第%d/%d页 %d条数据 使用了：%s", database, tableName, int64(pageIndex), int64(pageCount), len(results), sw.GetMillisecondsText())
		}
		// 导出的数据与查到的数据数量不一致
		if totalCount != realTotalCount {
			return lstBackupHistoryData, fmt.Errorf("%s.%s 导出的数据与查到的数据数量不一致", database, tableName)
		}
	}

	// 压缩
	cmd := fmt.Sprintf("gzip %s", backupRoot+fileName)
	code, result := exec.RunShellCommand(cmd, nil, path, false)
	if code != 0 {
		return lstBackupHistoryData, fmt.Errorf("压纹文件%s 时失败：%s", cmd, collections.NewList(result...).ToString(","))
	}
	fileName += ".gz"
	fileInfo, err := os.Stat(backupRoot + fileName)
	if err != nil {
		return lstBackupHistoryData, fmt.Errorf("获取备份文件信息:%s,失败： %s", fileName, err.Error())
	}

	lstBackupHistoryData.Add(BackupHistoryData{
		BackupId:  receiver.Id,
		Database:  database,
		FileName:  fileName,
		StoreType: receiver.StoreType,
		CreateAt:  dateTime.Now(),
		Size:      fileInfo.Size() / 1024,
	})
	return lstBackupHistoryData, nil
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

	var err error
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

	// 解压文件
	path := filepath.Dir(backupRoot + fileName)
	cmd := fmt.Sprintf("gzip -df %s", backupRoot+fileName)
	code, result := exec.RunShellCommand(cmd, nil, path, false)
	if code != 0 {
		return fmt.Errorf("解压文件：%s 时失败：%s", cmd, collections.NewList(result...).ToString(","))
	}
	fileName = fileName[:len(fileName)-3]
	defer file.Delete(backupRoot + fileName)

	// 恢复数据库
	switch receiver.BackupDataType {
	case eumBackupDataType.Mysql:
		err = db.RecoverMysql(receiver.Host, receiver.Port, receiver.Username, receiver.Password, database, backupRoot+fileName)
	case eumBackupDataType.Clickhouse:
		err = receiver.RecoverClickhouse(database, backupRoot+fileName)
	}
	return err
}

// 恢复clickhouse
func (receiver *DomainObject) RecoverClickhouse(database string, fileName string) error {
	// 读取表列表
	dbConnectionString := data.CreateConnectionString(receiver.BackupDataType.ToString(), receiver.Host, receiver.Port, database, receiver.Username, receiver.Password)
	dbContext := data.NewInternalContext(dbConnectionString)
	fSql, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("打开 %s 的SQL 文件失败: %v", fileName, err)
	}
	defer fSql.Close()

	// 逐行读取 SQL 文件
	scanner := bufio.NewScanner(fSql)
	var sqlBuilder strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		// 忽略空行和注释
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "--") {
			continue
		}

		// 拼接 SQL 语句
		sqlBuilder.WriteString(line + "\n")

		// 如果遇到分号，表示一个完整的 SQL 语句
		if strings.HasSuffix(line, ";") {
			sqlStatement := sqlBuilder.String()
			sqlBuilder.Reset()

			// 执行 SQL 语句
			if _, err := dbContext.ExecuteSql(sqlStatement); err != nil {
				return fmt.Errorf("执行 %s 的SQL 失败: %v", fileName, err)
			}
		}
	}

	// 检查文件读取错误
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取 %s SQL 的文件失败: %v", fileName, err)
	}
	return nil
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
