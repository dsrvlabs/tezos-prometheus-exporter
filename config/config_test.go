package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFile(t *testing.T) {
	fixtureFile := "./config_fixture.json"

	loader := NewLoader(fixtureFile)

	cfg, err := loader.Load()

	assert.Nil(t, err)
	assert.Equal(t, "http://localhost:8732", cfg.RPCEndpoint)
	assert.Equal(t, 9489, cfg.ExporterPort)
	assert.Equal(t, 1, cfg.UpdateIntervalSeconds)
	assert.Equal(t, "/mnt/tezos", cfg.DataDir)
}

func TestLoaderDefault(t *testing.T) {
	loader := NewLoader("")

	cfg, err := loader.Load()

	assert.Nil(t, err)
	assert.Equal(t, defaultRPCEndpoint, cfg.RPCEndpoint)
	assert.Equal(t, defaultServicePort, cfg.ExporterPort)
	assert.Equal(t, defaultUpdateInterval, cfg.UpdateIntervalSeconds)
	assert.Equal(t, defaultMountPath, cfg.DataDir)
}
