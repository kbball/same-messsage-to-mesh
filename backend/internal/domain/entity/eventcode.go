package entity

import "time"

// EventCode is a SAME/EAS event type code with its human-readable description.
type EventCode struct {
	Code        string
	Description string
	Category    string // Watch, Warning, Advisory, Statement, Test, etc.
	IsWarning   bool
	UpdatedAt   time.Time
}
