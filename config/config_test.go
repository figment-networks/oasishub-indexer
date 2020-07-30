package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromEnv(t *testing.T) {
	config := Config{}
	err := FromEnv(&config)

	assert.NoError(t, err)
	assert.Equal(t, modeDevelopment, config.AppEnv)
	assert.Equal(t, "0.0.0.0", config.ServerAddr)
	assert.Equal(t, int64(8081), config.ServerPort)
	assert.Equal(t, "@every 15m", config.IndexWorkerInterval)
	assert.Equal(t, int64(1), config.FirstBlockHeight)
	assert.Equal(t, false, config.Debug)
}

func TestListenAddr(t *testing.T) {
	config := Config{
		ServerAddr: "127.0.0.1",
		ServerPort: 5000,
	}
	assert.Equal(t, "127.0.0.1:5000", config.ListenAddr())
}

func TestValidate(t *testing.T) {
	config := Config{}
	assert.Equal(t, config.Validate(), errEndpointRequired)

	config.ProxyUrl = "endpoint"
	assert.Equal(t, config.Validate(), errDatabaseRequired)

	config.DatabaseDSN = "database"
	assert.NotEqual(t, config.Validate(), errDatabaseRequired)

	config.IndexWorkerInterval = ""
	assert.Equal(t, config.Validate(), errIndexWorkerIntervalRequired)
}
