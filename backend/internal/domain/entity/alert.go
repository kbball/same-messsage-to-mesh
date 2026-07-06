package entity

import "time"

// SAMEAlert is a decoded SAME/EAS weather alert message.
type SAMEAlert struct {
	ID         int64     `json:"id"`
	ReceivedAt time.Time `json:"received_at"`
	Originator string    `json:"originator"`  // WXR, EAS, CIV, PEP
	EventCode  string    `json:"event_code"`  // 3-letter SAME event code (e.g. TOR, SVR, RWT)
	FIPSCodes  []string  `json:"fips_codes"`  // affected area codes in "PSSCCC" format
	PurgeTime  string    `json:"purge_time"`  // HHMM duration until alert expires
	IssueTime  string    `json:"issue_time"`  // JJJHHMM (Julian day + HHMM)
	CallSign   string    `json:"call_sign"`   // NWS office call sign
	RawMessage string    `json:"raw_message"` // full decoded SAME header string
	Published  bool      `json:"published"`   // true once published to MQTT
}
