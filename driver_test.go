package driver

import (
	"strings"
	"testing"
	"time"
)

func TestEgoDriver_RegisterGrpcService(t *testing.T) {

	// etcd
	target := "etcd://localhost:2379"
	serviceName := "grpc.dtmserver"
	endpoint := "localhost:36790"
	driver := new(egoDriver)
	if err := driver.RegisterGrpcService(target, strings.Join([]string{endpoint, serviceName}, "/")); err != nil {
		t.Errorf("register consul fail err :%+v", err)
	}

	// nacos
	// target := "nacos://127.0.0.1:8848/dtmservice?namespaceId=public&timeoutMs=3000&notLoadCacheAtStart=true&logLevel=debug"
	// endpoint := "localhost:36790"
	// driver := new(zeroDriver)
	// if err := driver.RegisterGrpcService(target, endpoint); err != nil {
	// 	t.Errorf("register nacos fail err :%+v", err)
	// }

	time.Sleep(60 * time.Second)
}
