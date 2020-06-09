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
	assert.Equal(t, 8081, config.ServerPort)
	assert.Equal(t, "10s", config.IndexWorkerInterval)
	assert.Equal(t, "10m", config.CleanupInterval)
	assert.Equal(t, 1, config.FirstBlockHeight)
	assert.Equal(t, 1000, config.CleanupThreshold)
	assert.Equal(t, false, config.Debug)
}

func TestFromFile(t *testing.T) {
	config := Config{}
	assert.Error(t, FromFile("nonexist", &config), "no such file or directory")
	assert.NoError(t, FromFile("../test/fixtures/config.json", &config))
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

	config.IndexWorkerInterval = "10sec"
	assert.Equal(t, config.Validate(), errSyncIntervalInvalid)

	config.IndexWorkerInterval = "10s"
	assert.NotEqual(t, config.Validate(), errSyncIntervalInvalid)

	config.CleanupInterval = ""
	assert.Equal(t, config.Validate(), errCleanupIntervalRequired)

	config.CleanupInterval = "10sec"
	assert.Equal(t, config.Validate(), errCleanupIntervalInvalid)

	config.CleanupInterval = "10s"
	assert.NotEqual(t, config.Validate(), errCleanupIntervalInvalid)
}
