package eumK8SControllers

type Enum int

const (
	Deployment  Enum = iota // Deployment 无状态应用
	StatefulSet             // StatefulSet 有状态应用
	DaemonSet               // DaemonSet 所有节点都会运行一个实例
	Job                     // Job 一次性任务
	Cronjob                 // Cronjob 定时任务
)

// String 获取标签名称
func (eum Enum) String() string {
	switch eum {
	case Deployment:
		return "deployment"
	case StatefulSet:
		return "statefulSet"
	case DaemonSet:
		return "daemonSet"
	case Job:
		return "job"
	case Cronjob:
		return "cronjob"
	}
	return ""
}
