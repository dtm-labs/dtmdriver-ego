package driver

import (
	"testing"
	"time"
)

func TestEgoDriver_RegisterGrpcService(t *testing.T) {

	// etcd
	target := "etcd:///localhost:2379/dtmservice?userName=root&password=123456"
	// serviceName := "grpc.dtmserver"
	endpoint := "localhost:36790"
	driver := new(egoDriver)
	if err := driver.RegisterService(target, endpoint); err != nil {
		t.Errorf("register consul fail err :%+v", err)
	}

	time.Sleep(60 * time.Second)
}
