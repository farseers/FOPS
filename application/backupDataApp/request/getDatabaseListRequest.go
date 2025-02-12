package request

import (
	"fops/domain/_/eumBackupDataType"

	"github.com/farseer-go/webapi/check"
)

type GetDatabaseListRequest struct {
	BackupDataType eumBackupDataType.Enum // 备份数据类型
	Host           string                 // 主机
	Port           int                    // 端口
	Username       string                 // 用户名
	Password       string                 // 密码
}

func (receiver *GetDatabaseListRequest) Check() {
	// 主机
	check.IsTrue(len(receiver.Host) == 0, 403, "主机不能为空")
	// 端口
	check.IsTrue(receiver.Port < 1, 403, "端口不能为空")
	// 用户名
	check.IsTrue(len(receiver.Username) == 0, 403, "用户名不能为空")
}
