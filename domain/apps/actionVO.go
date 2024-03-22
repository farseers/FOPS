package apps

import (
	"fmt"
	"github.com/farseer-go/fs/configure"
	"github.com/farseer-go/fs/parse"
	"github.com/farseer-go/utils/http"
	"strings"
)

type stepVO struct {
	Name              string         // 名称
	ActionName        string         // action 名称
	ActionVer         string         // action 版本
	ActionDownloadUrl string         // action下载地址
	RepositoryName    string         // 仓库名称
	With              map[string]any // 参数
	Run               []string       // 运行脚本
}

func (receiver *stepVO) GetActionPath() string {
	return ActionsRoot + receiver.RepositoryName + "/" + receiver.ActionVer + "/" + receiver.ActionName
}

type ActionVO struct {
	Name   string // 工作流名称
	RunsOn string // 基础镜像系统
	Proxy  string // 代理
	Env    map[string]string
	Steps  []stepVO // 步骤
}

func LoadWorkflows(workflowsYmlPath string, appName string, gitName string) (ActionVO, error) {
	var (
		workflowsYmlContent string
		err                 error
		statusCode          int
	)

	// 支持读取失败时，尝试3次读取
	for i := 0; i < 3; i++ {
		// 通过http读取工作流定义的内容
		if workflowsYmlContent, statusCode, err = http.RequestProxy("GET", workflowsYmlPath, nil, nil, "", 2000, configure.GetString("Fops.GitAgent")); err == nil && statusCode == 200 {
			break
		}
	}

	if err != nil {
		return ActionVO{}, fmt.Errorf("读取WorkflowsYml错误：%s", err.Error())
	}

	if workflowsYmlContent == "" {
		return ActionVO{}, fmt.Errorf("WorkflowsYml没有定义。")
	}

	if statusCode != 200 {
		return ActionVO{}, fmt.Errorf("读取WorkflowsYml失败：%d", statusCode)
	}

	// 替换项目名称
	workflowsYmlContent = strings.ReplaceAll(workflowsYmlContent, "${app_name}", appName)
	workflowsYmlContent = strings.ReplaceAll(workflowsYmlContent, "${git_name}", gitName)

	workflowsYml := configure.NewYamlConfig("")
	err = workflowsYml.LoadContent([]byte(workflowsYmlContent))
	if err != nil {
		return ActionVO{}, fmt.Errorf("读取WorkflowsYml错误：%s", err.Error())
	}

	name, _ := workflowsYml.Get("name")
	proxy, _ := workflowsYml.Get("jobs.build.proxy")
	sysImage, _ := workflowsYml.Get("jobs.build.runs-on")
	env, _ := workflowsYml.GetSubNodes("jobs.build.env")

	act := ActionVO{
		Name:   strings.TrimSpace(parse.ToString(name)),
		Proxy:  strings.TrimSpace(parse.ToString(proxy)),
		RunsOn: strings.TrimSpace(parse.ToString(sysImage)),
		Env:    make(map[string]string),
	}
	for k, v := range env {
		act.Env[k] = parse.ToString(v)
	}

	// 移除前缀//
	if index := strings.Index(act.Proxy, "//"); index > -1 {
		act.Proxy = act.Proxy[index+2:]
	}

	// 运行steps
	if steps, existsSteps := workflowsYml.Get("jobs.build.steps"); existsSteps {
		stepsLength := len(steps.([]any))
		for i := 0; i < stepsLength; i++ {
			step := stepVO{
				With: make(map[string]any),
			}
			curSteps := fmt.Sprintf("jobs.build.steps[%d].", i)
			// steps.name
			if curStepsName, b := workflowsYml.Get(curSteps + "name"); b {
				step.Name = parse.ToString(curStepsName)
			}

			// steps.uses
			if curStepsUses, b := workflowsYml.Get(curSteps + "uses"); b {
				step.ActionName = strings.Split(curStepsUses.(string), "@")[0]
				step.ActionVer = strings.Split(curStepsUses.(string), "@")[1]
				// 第三方action eg. farseers/FOPS-Actions/checkout@v1
				if strings.Contains(step.ActionName, "/") {
					step.RepositoryName = step.ActionName[0:strings.LastIndex(step.ActionName, "/")] // eg. farseers/FOPS-Actions
					step.ActionName = step.ActionName[strings.LastIndex(step.ActionName, "/")+1:]    // eg. checkout
				} else {
					step.RepositoryName = "farseers/FOPS-Actions"
				}

				// https://github.com/farseer-go/fsctl/releases/download/v0.13.1/fsctl.Darwin.x86_64
				step.ActionDownloadUrl = fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", step.RepositoryName, step.ActionVer, step.ActionName)
			}
			// steps.with
			if curStepsWith, b := workflowsYml.GetSubNodes(curSteps + "with"); b {
				step.With = curStepsWith
			}

			// steps.run
			if curStepsRun, b := workflowsYml.Get(curSteps + "run"); b {
				for _, run := range curStepsRun.([]any) {
					step.Run = append(step.Run, parse.ToString(run))
				}
			}
			act.Steps = append(act.Steps, step)
		}
	}
	return act, nil
}
