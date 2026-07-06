package entity

import "time"

// FIPSCode maps a SAME FIPS state+county code pair to its human-readable names.
type FIPSCode struct {
	StateCode  string // 2-digit FIPS state code (e.g. "13" for Georgia)
	CountyCode string // 3-digit FIPS county code (e.g. "121" for Fulton County)
	StateName  string
	CountyName string
	UpdatedAt  time.Time
}

// FIPS returns the combined 5-character FIPS code used in SAME messages.
func (f FIPSCode) FIPS() string {
	return f.StateCode + f.CountyCode
}
