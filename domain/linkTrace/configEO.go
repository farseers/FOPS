package linkTrace

var Config ConfigEO

// ConfigEO 配置 ./farseer.yaml.FOPS.LinkTrace
type ConfigEO struct {
	Driver     string // 数据库类型
	ConnString string // 连接字符串
}
