package service

import (
	"context"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

// AlertPublisher sends a decoded alert to an external transport (e.g. MQTT).
type AlertPublisher interface {
	Publish(ctx context.Context, alert entity.SAMEAlert, message string) error
}

// SAMEDecoder reads from the SDR hardware pipeline and emits decoded alerts.
// Implementations wrap the rtl_fm | multimon-ng process.
type SAMEDecoder interface {
	// Start begins the decoding pipeline. Decoded alerts are sent to ch.
	// The caller is responsible for draining ch and calling Stop.
	Start(ctx context.Context, ch chan<- entity.SAMEAlert) error
	Stop()
}
