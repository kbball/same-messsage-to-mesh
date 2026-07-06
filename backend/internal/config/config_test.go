package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_Defaults(t *testing.T) {
	// Clear any vars that .env may export so defaults are exercised.
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_PORT", "")
	t.Setenv("SERVER_PORT", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("DB_NAME", "")
	t.Setenv("SDR_FREQUENCY", "")
	t.Setenv("MQTT_ENABLED", "")
	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, 8080, cfg.ServerPort)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.Equal(t, 5432, cfg.DB.Port)
	assert.Equal(t, "same_mesh", cfg.DB.Name)
	assert.Equal(t, false, cfg.MQTT.Enabled)
	assert.Equal(t, int64(162550000), cfg.SDR.Frequency)
}

func TestLoad_MissingPassword(t *testing.T) {
	t.Setenv("DB_PASSWORD", "")
	_, err := Load()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DB_PASSWORD")
}

func TestDBConfig_DSN(t *testing.T) {
	cfg := DBConfig{
		Host:     "localhost",
		Port:     5432,
		Name:     "same_mesh",
		User:     "postgres",
		Password: "secret",
		SSLMode:  "disable",
	}
	dsn := cfg.DSN()
	assert.Contains(t, dsn, "host=localhost")
	assert.Contains(t, dsn, "port=5432")
	assert.Contains(t, dsn, "dbname=same_mesh")
}
