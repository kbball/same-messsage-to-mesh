package entity

import "time"

// SAMEAlert is a decoded SAME/EAS weather alert message.
type SAMEAlert struct {
	ID         int64
	ReceivedAt time.Time
	Originator string    // WXR, EAS, CIV, PEP
	EventCode  string    // 3-letter SAME event code (e.g. TOR, SVR, RWT)
	FIPSCodes  []string  // affected area codes in "PSSCCC" format
	PurgeTime  string    // HHMM duration until alert expires
	IssueTime  string    // JJJHHMM (Julian day + HHMM)
	CallSign   string    // NWS office call sign
	RawMessage string    // full decoded SAME header string
	Published  bool      // true once published to MQTT
}
