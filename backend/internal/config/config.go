package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	ServerPort int
	LogLevel   string
	DB         DBConfig
	MQTT       MQTTConfig
	SDR        SDRConfig
}

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		c.Host, c.Port, c.Name, c.User, c.Password, c.SSLMode,
	)
}

type MQTTConfig struct {
	Enabled      bool
	Host         string
	Port         int
	PublishTopic string
}

type SDRConfig struct {
	DevicePath string
	Frequency  int64
}

func Load() (*Config, error) {
	serverPort, err := envInt("SERVER_PORT", 8080)
	if err != nil {
		return nil, fmt.Errorf("SERVER_PORT: %w", err)
	}

	dbPort, err := envInt("DB_PORT", 5432)
	if err != nil {
		return nil, fmt.Errorf("DB_PORT: %w", err)
	}

	mqttEnabled, err := envBool("MQTT_ENABLED", false)
	if err != nil {
		return nil, fmt.Errorf("MQTT_ENABLED: %w", err)
	}

	mqttPort, err := envInt("MQTT_PORT", 1883)
	if err != nil {
		return nil, fmt.Errorf("MQTT_PORT: %w", err)
	}

	sdrFrequency, err := envInt64("SDR_FREQUENCY", 162550000)
	if err != nil {
		return nil, fmt.Errorf("SDR_FREQUENCY: %w", err)
	}

	cfg := &Config{
		ServerPort: serverPort,
		LogLevel:   envStr("LOG_LEVEL", "info"),
		DB: DBConfig{
			Host:     envStr("DB_HOST", "localhost"),
			Port:     dbPort,
			Name:     envStr("DB_NAME", "same_mesh"),
			User:     envStr("DB_USER", "postgres"),
			Password: envStr("DB_PASSWORD", ""),
			SSLMode:  envStr("DB_SSL_MODE", "disable"),
		},
		MQTT: MQTTConfig{
			Enabled:      mqttEnabled,
			Host:         envStr("MQTT_HOST", "localhost"),
			Port:         mqttPort,
			PublishTopic: envStr("MQTT_PUBLISH_TOPIC", "same/alerts"),
		},
		SDR: SDRConfig{
			DevicePath: envStr("SDR_DEVICE_PATH", "/dev/bus/usb"),
			Frequency:  sdrFrequency,
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	var missing []string
	if c.DB.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}
	return nil
}

func envStr(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func envInt(key string, defaultVal int) (int, error) {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", v)
	}
	return n, nil
}

func envInt64(key string, defaultVal int64) (int64, error) {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal, nil
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", v)
	}
	return n, nil
}

func envBool(key string, defaultVal bool) (bool, error) {
	v := os.Getenv(key)
	if v == "" {
		return defaultVal, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("invalid boolean %q", v)
	}
	return b, nil
}
