package entity

import "time"

// AlertFilter is the user's configuration for which SAME alerts to act on.
// Empty slices mean "all" for that dimension.
type AlertFilter struct {
	StateCodes []string  `json:"state_codes"` // FIPS state codes to include; empty = all states
	FIPSCodes  []string  `json:"fips_codes"`  // specific county FIPS codes; empty = all counties in selected states
	EventCodes []string  `json:"event_codes"` // SAME event codes to include; empty = all event types
	UpdatedAt  time.Time `json:"updated_at"`
}

// SDRDeviceConfig is the persisted SDR hardware configuration.
type SDRDeviceConfig struct {
	DevicePath string    `json:"device_path"`
	Frequency  int64     `json:"frequency"` // Hz (e.g. 162550000 for 162.550 MHz)
	UpdatedAt  time.Time `json:"updated_at"`
}

// MQTTConfig is the persisted MQTT broker configuration.
type MQTTConfig struct {
	Enabled      bool      `json:"enabled"`
	Host         string    `json:"host"`
	Port         int       `json:"port"`
	PublishTopic string    `json:"publish_topic"`
	UpdatedAt    time.Time `json:"updated_at"`
}
