package driver

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/dtm-labs/dtmdriver"
	"github.com/gotomicro/ego/core/constant"
	"github.com/gotomicro/ego/server"
	"github.com/vicnoah/dtmdriver-ego/plugin/egodriver"
)

const (
	DriverName = "dtm-driver-ego"
	KindEtcd   = "etcd" // etcd
	KindK8S    = "k8s"  // k8s
)

type (
	egoDriver struct{}
)

func (e *egoDriver) GetName() string {
	return DriverName
}

func (e *egoDriver) RegisterGrpcResolver() {
	egodriver.Register()
}

func (e *egoDriver) RegisterGrpcService(target string, endpoint string) error {
	if target == "" { // empty target, no action
		return nil
	}
	target = strings.Replace(target, "///", "//", -1)
	u, err := url.Parse(target)
	if err != nil {
		return err
	}

	info := server.ApplyOptions(
		server.WithScheme("grpc"),
		server.WithAddress(endpoint),
		server.WithKind(constant.ServiceProvider),
		server.WithName(u.Host),
	)

	return egodriver.RegisterMap[u.Scheme].RegisterService(context.Background(), &info)
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
	index := strings.IndexByte(u.Path[1:], '/') + 1
	return u.Scheme + "://" + u.Host, u.Path[index:], nil
}

func init() {
	dtmdriver.Register(&egoDriver{})
}
