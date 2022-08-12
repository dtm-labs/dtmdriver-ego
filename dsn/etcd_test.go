package dsn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEtcdDsnParser_ParseDSN(t *testing.T) {
	dsn := "etcd://localhost:2357,192.168.0.1:2357/dtmservice"
	var parser EtcdDSNParser
	cfg, err := parser.ParseDSN(dsn)
	assert.NoError(t, err)
	assert.Equal(t, "etcd", cfg.Scheme)
	assert.Equal(t, "dtmservice", cfg.ServiceName)
	fmt.Println(cfg)
}
