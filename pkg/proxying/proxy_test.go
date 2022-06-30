package proxying

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wamphlett/ledblinky-proxy/config"
)

func TestConfiguresProxyCorrectly(t *testing.T) {
	proxy := &Proxy{}

	// create a config for multiple publishers
	cfg := &config.ReceiversConfig{
		Executables: []string{"first.exe", "second.exe"},
		Webhooks:    []string{"localhost:3000", "remote:5000"},
	}

	proxy.ConfigurePublishers(cfg)
	assert.Len(t, proxy.publishers, 4)
}

func TestConfiguresProxyWithNoReceivers(t *testing.T) {
	proxy := &Proxy{}

	// create an empty config
	cfg := &config.ReceiversConfig{}

	proxy.ConfigurePublishers(cfg)
	assert.Len(t, proxy.publishers, 0)
}

func TestConfiguresProxyWithInvalidReceivers(t *testing.T) {
	proxy := &Proxy{}

	// create an empty config
	cfg := &config.ReceiversConfig{
		Executables: []string{""},
	}

	proxy.ConfigurePublishers(cfg)
	assert.Len(t, proxy.publishers, 0)
}
