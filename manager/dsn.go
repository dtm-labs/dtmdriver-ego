package manager

import "github.com/gotomicro/ego/core/eregistry"

// DSN ...
type DSN struct {
	Scheme      string
	ServiceName string
	Registry    eregistry.Registry
}

// DSNParser ...
type DSNParser interface {
	ParseDSN(string) (cfg *DSN, err error)
	Scheme() string
}
