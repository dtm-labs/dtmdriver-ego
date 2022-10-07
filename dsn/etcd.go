package dsn

import (
	"bytes"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dtm-labs/dtmdriver-ego/manager"
	"github.com/gotomicro/ego-component/eetcd"
	"github.com/gotomicro/ego-component/eetcd/registry"
	"github.com/gotomicro/ego/core/econf"
)

var (
	_ manager.DSNParser = (*EtcdDSNParser)(nil)
)

// EtcdDSNParser ...
type EtcdDSNParser struct {
}

func init() {
	manager.Register(&EtcdDSNParser{})
}

// Scheme ...
func (p *EtcdDSNParser) Scheme() string {
	return "etcd"
}

// ParseDSN ...
func (p *EtcdDSNParser) ParseDSN(dsn string) (cfg *manager.DSN, err error) {
	u, err := url.Parse(dsn)
	if err != nil || u.Host == "" {
		return
	}
	cfg = new(manager.DSN)
	cfg.Scheme = p.Scheme()
	cfg.ServiceName = strings.TrimPrefix(u.Path, "/")

	query := u.Query()
	conf := eetcd.DefaultConfig()

	conf.Addrs = strings.Split(u.Host, ",")
	if query.Get("certFile") != "" {
		conf.CertFile = query.Get("certFile")
	}
	if query.Get("keyFile") != "" {
		conf.KeyFile = query.Get("keyFile")
	}
	if query.Get("caCert") != "" {
		conf.CaCert = query.Get("caCert")
	}
	if query.Get("userName") != "" {
		conf.UserName = query.Get("userName")
	}
	if query.Get("connectTimeout") != "" {
		du, _ := time.ParseDuration(query.Get("connectTimeout"))
		conf.ConnectTimeout = du
	}
	if query.Get("autoSyncInterval") != "" {
		du, _ := time.ParseDuration(query.Get("autoSyncInterval"))
		conf.AutoSyncInterval = du
	}
	if query.Get("enableBasicAuth") != "" {
		b, _ := strconv.ParseBool(query.Get("enableBasicAuth"))
		conf.EnableBasicAuth = b
	}
	if query.Get("enableSecure") != "" {
		b, _ := strconv.ParseBool(query.Get("enableSecure"))
		conf.EnableSecure = b
	}
	if query.Get("enableBlock") != "" {
		b, _ := strconv.ParseBool(query.Get("enableBlock"))
		conf.EnableBlock = b
	}
	if query.Get("enableFailOnNonTempDialError") != "" {
		b, _ := strconv.ParseBool(query.Get("enableFailOnNonTempDialError"))
		conf.EnableFailOnNonTempDialError = b
	}
	var buf = new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("[%s]\n", u.Scheme))
	toml.NewEncoder(buf).Encode(&conf)
	buf.WriteString(fmt.Sprintf("\n[registry]\nscheme = \"%s\"", u.Scheme))
	econf.LoadFromReader(buf, toml.Unmarshal)

	cfg.Registry = registry.Load("registry").Build(registry.WithClientEtcd(eetcd.Load(u.Scheme).Build()))

	return
}
