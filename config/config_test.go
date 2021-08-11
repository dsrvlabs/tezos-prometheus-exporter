package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfig(t *testing.T) {

	tests := []struct {
		Desc      string
		MockArgs  []string
		ExpectCfg Config
	}{
		{
			Desc: "Success parsing config",
			MockArgs: []string{
				"-rpc-endpoint", "hello",
				"-prometheus-port", "1234",
				"-fetch-interval", "10",
				"-mount-path", "/mnt/tezos",
			},
			ExpectCfg: Config{
				RPCEndpoint:           "hello",
				ServicePort:           1234,
				UpdateIntervalSeconds: 10,
				MountPath:             "/mnt/tezos",
			},
		},

		{
			Desc:     "Omit Value",
			MockArgs: []string{},
			ExpectCfg: Config{
				RPCEndpoint:           defaultRPCEndpoint,
				ServicePort:           defaultServicePort,
				UpdateIntervalSeconds: defaultUpdateInterval,
				MountPath:             defaultMountPath,
			},
		},
	}

	for _, test := range tests {
		args := []string{"./app"}
		args = append(args, test.MockArgs...)

		// Mock
		os.Args = args

		// Test
		// Intentionally call parserConfig function instread of GetConfig
		// because GetConfig function uses singleton pattern so the config instance is not updated.
		cfg := parseConfig()

		// Assert
		assert.Equal(t, test.ExpectCfg.RPCEndpoint, cfg.RPCEndpoint)
		assert.Equal(t, test.ExpectCfg.ServicePort, cfg.ServicePort)
		assert.Equal(t, test.ExpectCfg.UpdateIntervalSeconds, cfg.UpdateIntervalSeconds)
		assert.Equal(t, test.ExpectCfg.MountPath, cfg.MountPath)
	}
}
