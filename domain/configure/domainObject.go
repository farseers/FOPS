package configure

// DomainObject 配置中心
type DomainObject struct {
	AppName string // 应用名称
	Key     string // 配置KEY
	Ver     int    // 版本
	Value   string // 配置VALUE
}

func (receiver DomainObject) IsNil() bool {
	return receiver.AppName == ""
}
