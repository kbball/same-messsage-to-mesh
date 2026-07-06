package entity

import "time"

// FIPSCode maps a SAME FIPS state+county code pair to its human-readable names.
type FIPSCode struct {
	StateCode  string    `json:"state_code"`  // 2-digit FIPS state code (e.g. "13" for Georgia)
	CountyCode string    `json:"county_code"` // 3-digit FIPS county code (e.g. "121" for Fulton County)
	StateName  string    `json:"state_name"`
	CountyName string    `json:"county_name"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// FIPS returns the combined 5-character FIPS code used in SAME messages.
func (f FIPSCode) FIPS() string {
	return f.StateCode + f.CountyCode
}
