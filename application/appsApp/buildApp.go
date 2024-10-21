// @area /apps/
package appsApp

import (
	"fmt"
	"fops/application/appsApp/response"
	"fops/domain/apps"
	"fops/domain/cluster"
	"regexp"
	"strings"

	"github.com/farseer-go/collections"
	"github.com/farseer-go/fs/core"
	"github.com/farseer-go/fs/dateTime"
	"github.com/farseer-go/fs/exception"
	"github.com/farseer-go/mapper"
	"github.com/farseer-go/webapi/action"
)

// BuildAdd 添加构建
// @post build/add
// @filter application.Jwt
func BuildAdd(appName string, workflowsName string, appsRepository apps.Repository, clusterRepository cluster.Repository, dockerDevice apps.IDockerDevice) {
	appDO := appsRepository.ToEntity(appName)
	exception.ThrowWebExceptionfBool(appDO.IsNil(), 403, "应用不存在")
	exception.ThrowWebExceptionfBool(workflowsName == "", 403, "工作流名称未设置")

	buildNumber := appsRepository.GetBuildNumber(appName) + 1
	buildDO := apps.BuildEO{
		BuildServerId: core.AppId,
		BuildNumber:   buildNumber,
		CreateAt:      dateTime.Now(),
		FinishAt:      dateTime.Now(),
		Env:           apps.EnvVO{},
		AppName:       appName,
		WorkflowsName: workflowsName,
	}
	err := appsRepository.AddBuild(buildDO)
	exception.ThrowWebExceptionError(403, err)
}

// BuildList 构建列表
// @post build/list
// @filter application.Jwt
func BuildList(appName string, pageSize int, pageIndex int, appsRepository apps.Repository) collections.PageList[response.BuildListResponse] {
	if pageSize < 1 {
		pageSize = 20
	}
	if pageIndex < 1 {
		pageIndex = 1
	}
	lst := appsRepository.ToBuildList(appName, pageSize, pageIndex)
	return mapper.ToPageList[response.BuildListResponse](lst, func(r *response.BuildListResponse, a any) {
		buildEO := a.(apps.BuildEO)
		r.CreateAt = buildEO.CreateAt.ToString("MM-dd HH:mm:ss")
		r.FinishAt = buildEO.FinishAt.ToString("MM-dd HH:mm:ss")
	})
}

// 语法高亮
var chineseTips = collections.NewList(
	"环境变量：", "前置检查。", "先删除之前编译的目标文件。", "自动创建目录。", "前置检查通过。", "已经是最新的。", "拉取完成。", "登陆镜像仓库。", "镜像仓库登陆成功。",
	"开始镜像打包。", "镜像打包完成。", "开始上传镜像。", "镜像上传完成。", "开始更新K8S POD的镜像版本。", "更新镜像版本完成。")

// （黄色）
var cmdPrefix = collections.NewList(
	"git -C ",
	"docker login ",
	"docker build -t",
	"docker push ",
	"kubectl set image ",
	"Successfully built ",
	"Successfully tagged ",
	"成功执行。",
	"构建完成。")

var cmdResultTips = collections.NewList(
	"Login Succeeded",
	"The push refers to repository ",
	"Determining projects to restore...")

var errorTips = collections.NewList(
	"Exception ",
	"Failed ",
	"ERROR:",
	"error ",
	"镜像打包出错了。",
	"镜像仓库登陆失败。",
	"执行失败，退出构建。",
	"K8S更新镜像失败。",
	"Unable to connect",
	"Cannot connect")

//var globalStr []string
//var id int64

// View 构建日志
// @get build/view-{buildId}
func View(buildId int64) action.IResult {
	//if id != buildId {
	//	globalStr = []string{}
	//}
	//id = buildId
	logQueue := apps.LogQueue{
		BuildId: buildId,
	}
	logContent := logQueue.View()
	for i := 0; i < len(logContent); i++ {
		if len(logContent[i]) < 20 {
			continue
		}
		dateTimePart := logContent[i][:19] // 提取日期时间部分
		logPart := logContent[i][20:]      // 提取日志内容部分
		dateTimePart = fmt.Sprintf("<span style=\"color:#9caf62\">%s</span>", dateTimePart)

		if chineseTips.Contains(logPart) {
			logPart = fmt.Sprintf("<span style=\"color:#cfbbfc\">%s</span>", logPart)
		} else if cmdResultTips.ContainsAny(logPart) || strings.HasPrefix(logPart, "The push refers to repository ") || regexp.MustCompile(`\w+\.apps/\w+ image updated`).MatchString(logPart) {
			logPart = fmt.Sprintf("<span style=\"color:#fff\">%s</span>", logPart)
		} else if cmdPrefix.ContainsAny(logPart) {
			logPart = fmt.Sprintf("<span style=\"color:#ffe127\">%s</span>", logPart)
		} else if errorTips.ContainsAny(strings.ToLower(logPart)) {
			logPart = fmt.Sprintf("<span style=\"color:#ff5b5b\">%s</span>", logPart)
		} else {
			dockerLogMatch := regexp.MustCompile(`#\d+ \[.+ \d+/\d+\] (?P<cmd>.+)`).FindStringSubmatch(logPart)
			if dockerLogMatch == nil {
				dockerLogMatch = regexp.MustCompile(`#\d+ \d+\.\d+ (?P<cmd>.+)`).FindStringSubmatch(logPart)
			}
			if dockerLogMatch == nil {
				dockerLogMatch = regexp.MustCompile(`Step \d+/\d+ :(?P<cmd>.+)`).FindStringSubmatch(logPart)
			}
			if dockerLogMatch != nil {
				cmd := dockerLogMatch[1]
				logPart = strings.Replace(logPart, cmd, fmt.Sprintf("<span style=\"color:#38e4c6\">%s</span>", cmd), 1)
				logPart = fmt.Sprintf("<span style=\"color:#a6a49a\">%s</span>", logPart)
			} else {
				logPart = fmt.Sprintf("<span style=\"color:#a6a49a\">%s</span>", logPart)
			}
		}
		logContent[i] = dateTimePart + " " + logPart
	}
	return action.Content(strings.Join(logContent, "\n"))
}

// Stop 停止构建
// @post build/stop
// @filter application.Jwt
func Stop(buildId int64, appsRepository apps.Repository) {
	var buildEO apps.BuildEO
	if buildId > 0 {
		buildEO = appsRepository.ToBuildEntity(buildId)
	} else {
		// 找到最后一个正在building的构建任务
		buildEO = appsRepository.GetLastBuilding()
	}
	buildEO.SetCancel()
}
