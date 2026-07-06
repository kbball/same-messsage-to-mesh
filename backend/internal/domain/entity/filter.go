package entity

import "time"

// AlertFilter is the user's configuration for which SAME alerts to act on.
// Empty slices mean "all" for that dimension.
type AlertFilter struct {
	StateCodes []string // FIPS state codes to include; empty = all states
	FIPSCodes  []string // specific county FIPS codes; empty = all counties in selected states
	EventCodes []string // SAME event codes to include; empty = all event types
	UpdatedAt  time.Time
}

// SDRDeviceConfig is the persisted SDR hardware configuration.
type SDRDeviceConfig struct {
	DevicePath string
	Frequency  int64 // Hz (e.g. 162550000 for 162.550 MHz)
	UpdatedAt  time.Time
}
