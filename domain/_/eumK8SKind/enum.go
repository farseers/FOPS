package eumK8SKind

type Enum int

const (
	Controllers Enum = iota
	Ingress
	Service
	Config
)
