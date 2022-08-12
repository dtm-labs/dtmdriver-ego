package dsn

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestK8sDsnParser_ParseDSN(t *testing.T) {
	dsn := "k8s://dtmservice:36790/?namespaces=default"
	var parser EtcdDSNParser
	cfg, err := parser.ParseDSN(dsn)
	assert.NoError(t, err)
	assert.Equal(t, "k8s", cfg.Scheme)
	assert.Equal(t, "dtmservice", cfg.ServiceName)
	fmt.Println(cfg)
}
