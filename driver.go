package driver

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/dtm-labs/dtmdriver"
	"github.com/gotomicro/ego/client/egrpc/resolver"
	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/core/eregistry"
	"github.com/gotomicro/ego/server"
	_ "github.com/vicnoah/dtmdriver-ego/dsn"
	"github.com/vicnoah/dtmdriver-ego/manager"
)

const (
	DriverName = "dtm-driver-ego"
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

	cfg, err := manager.Parse(target)
	if err != nil {
		return err
	}

	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(endpoint),
		server.WithKind(constant.ServiceProvider),
		server.WithName(cfg.ServiceName),
	)

	reg := cfg.Registry
	registerMap[cfg.Scheme] = reg

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
