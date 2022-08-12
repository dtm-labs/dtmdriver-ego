package driver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "github.com/vicnoah/dtmdriver-ego/dsn"
	"github.com/vicnoah/dtmdriver-ego/manager"
)

func TestEtcdDsnParser_ParseDSN(t *testing.T) {
	dsn := "etcd://localhost:2357,192.168.0.1:2357/dtmservice?scheme=etcd&caCert=/vector/cert.pem"
	cfg, err := manager.Parse(dsn)
	assert.NoError(t, err)
	assert.Equal(t, "etcd", cfg.Scheme)
	assert.Equal(t, "dtmservice", cfg.ServiceName)
	fmt.Println(cfg)
}
