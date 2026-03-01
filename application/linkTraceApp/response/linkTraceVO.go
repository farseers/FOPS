package response

import "github.com/farseer-go/fs/trace"

// "245,108,108,0.4"
var RgbaList = []string{"95,184,120,0.4", "65,105,225,0.4", "219,112,147,0.4", "128,0,128,0.4", "153,50,204,0.4", "123,104,238,0.4", "119,136,153,0.4", "70,130,180,0.4", "0,139,139,0.4", "34,139,34,0.4", "128,128,0,0.4", "238,232,170,0.4", "218,165,32,0.4", "255,165,0,0.4", "255,140,0,0.4", "210,105,30,0.4"}

type LinkTraceVO struct {
	Rgba        string  // 背景颜色
	AppId       string  // 应用实例ID
	AppIp       string  // 应用IP
	AppName     string  // 应用名称
	StartTs     float64 // 计算服务的相对调用时间，第一个服务从0us计算
	StartRate   float64 // 计算服务的相对调用时间的百分比
	UseTs       float64 // 使用时间（微秒）
	UseRate     float64 // 使用时间占比总用时的百分比
	UseDesc     string  // 使用时间（描述）
	Caption     string  // 标题
	Desc        string  // tips描述
	CopyContent string  // 用于前端复制时的内容
	Exception   *trace.ExceptionStack
}
