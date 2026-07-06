package entity

import "time"

// EventCode is a SAME/EAS event type code with its human-readable description.
type EventCode struct {
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Category    string    `json:"category"` // Watch, Warning, Advisory, Statement, Test, etc.
	IsWarning   bool      `json:"is_warning"`
	UpdatedAt   time.Time `json:"updated_at"`
}
