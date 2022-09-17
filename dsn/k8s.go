package dsn

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/gotomicro/ego-component/ek8s"
	"github.com/gotomicro/ego-component/ek8s/registry"
	"github.com/gotomicro/ego/core/econf"
	"github.com/vicnoah/dtmdriver-ego/manager"
)

var (
	_ manager.DSNParser = (*K8sDSNParser)(nil)
)

// K8sDSNParser ...
type K8sDSNParser struct {
}

func init() {
	manager.Register(&K8sDSNParser{})
}

// Scheme ...
func (p *K8sDSNParser) Scheme() string {
	return "k8s"
}

// ParseDSN ...
func (p *K8sDSNParser) ParseDSN(dsn string) (cfg *manager.DSN, err error) {
	u, err := url.Parse(dsn)
	if err != nil || u.Host == "" {
		return
	}
	cfg = new(manager.DSN)
	cfg.Scheme = p.Scheme()
	cfg.ServiceName = u.Host

	query := u.Query()
	conf := ek8s.DefaultConfig()

	conf.Addr = u.Host
	if query.Get("debug") != "" {
		b, _ := strconv.ParseBool(query.Get("debug"))
		conf.Debug = b
	}
	if query.Get("token") != "" {
		conf.Token = query.Get("token")
	}
	if query.Get("namespaces") != "" {
		conf.Namespaces = strings.Split(query.Get("namespaces"), ",")
	}
	if query.Get("deploymentPrefix") != "" {
		conf.DeploymentPrefix = query.Get("deploymentPrefix")
	}
	if query.Get("tlsClientConfigInsecure") != "" {
		b, _ := strconv.ParseBool(query.Get("tlsClientConfigInsecure"))
		conf.TLSClientConfigInsecure = b
	}
	var buf = new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("[%s]\n", u.Scheme))
	toml.NewEncoder(buf).Encode(&conf)
	buf.WriteString(fmt.Sprintf("\n[registry]\nscheme = \"%s\"", u.Scheme))
	econf.LoadFromReader(buf, toml.Unmarshal)

	cfg.Registry = registry.Load("registry").Build(registry.WithClient(ek8s.Load(u.Scheme).Build()))

	return
}
