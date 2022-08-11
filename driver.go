package driver

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dtm-labs/dtmdriver"
	"github.com/gotomicro/ego-component/eetcd"
	"github.com/gotomicro/ego-component/eetcd/registry"
	"github.com/gotomicro/ego-component/ek8s"
	k8sregistry "github.com/gotomicro/ego-component/ek8s/registry"
	"github.com/gotomicro/ego/client/egrpc/resolver"
	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/econf"
	"github.com/gotomicro/ego/core/eregistry"
	"github.com/gotomicro/ego/server"
)

const (
	DriverName = "dtm-driver-ego"
	KindEtcd   = "etcd" // etcd
	KindK8S    = "k8s"  // k8s
)

type (
	egoDriver struct{}
)

var registerMap = make(map[string]eregistry.Registry, 0)

func (e *egoDriver) GetName() string {
	return DriverName
}

func (e *egoDriver) RegisterAddrResolver() {
	for name, reg := range registerMap {
		resolver.Register(name, reg)
	}
}

func (e *egoDriver) RegisterService(target string, endpoint string) error {
	if target == "" { // empty target, no action
		return nil
	}
	target = strings.Replace(target, "///", "//", -1)
	u, err := url.Parse(target)
	if err != nil {
		return err
	}

	query, _ := url.ParseQuery(u.RawQuery)

	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(endpoint),
		server.WithKind(constant.ServiceProvider),
		server.WithName(u.Host),
	)

	var reg eregistry.Registry

	switch u.Scheme {
	case KindEtcd:
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

		reg = registry.Load("registry").Build(registry.WithClientEtcd(eetcd.Load(u.Scheme).Build()))
		registerMap[u.Scheme] = reg
	case KindK8S:
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

		reg = k8sregistry.Load("registry").Build(k8sregistry.WithClient(ek8s.Load(u.Scheme).Build()))
		registerMap[u.Scheme] = reg
	}

	if reg != nil {
		return errors.New("register error")
	}
	return reg.RegisterService(context.Background(), &info)
}

// etcd:///<服务名称>/api.hello/TransOut
// k8s:///<服务名称>/api.hello/TransOut
func (e *egoDriver) ParseServerMethod(uri string) (server string, method string, err error) {
	// direct 直连服务
	if !strings.Contains(uri, "///") {
		sep := strings.IndexByte(uri, '/')
		if sep == -1 {
			return "", "", fmt.Errorf("bad url: '%s'. no '/' found", uri)
		}
		return uri[:sep], uri[sep:], nil

	}

	uri = strings.Replace(uri, "///", "//", -1)

	u, err := url.Parse(uri)
	if err != nil {
		return "", "", nil
	}
	return u.Scheme + ":///" + u.Host, u.Path, nil
}

func init() {
	dtmdriver.Register(&egoDriver{})
}
