package egodriver

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/gotomicro/ego"
	"github.com/gotomicro/ego-component/eetcd"
	"github.com/gotomicro/ego-component/eetcd/registry"
	"github.com/gotomicro/ego/core/econf"
)

func init() {
	ego.New().Invoker(func() error {
		var buf = new(bytes.Buffer)
		buf.WriteString(`[etcd]
addrs=["127.0.0.1:2379"]
connectTimeout = "1s"
secure = false

[registry]
scheme = "etcd" # grpc resolver默认scheme为"etcd"，你可以自行修改`)

		return econf.LoadFromReader(buf, toml.Unmarshal)
	}, func() error {
		// 注册etcd节点
		com := registry.Load("registry").Build(registry.WithClientEtcd(eetcd.Load("etcd").Build()))
		RegisterMap["etcd"] = com
		return nil
	})
}
