package egodriver

import (
	"github.com/gotomicro/ego/client/egrpc/resolver"
	"github.com/gotomicro/ego/core/eregistry"
)

var RegisterMap = make(map[string]eregistry.Registry, 0)

// Register 注册解析器
func Register() {
	for name, reg := range RegisterMap {
		resolver.Register(name, reg)
	}
}
