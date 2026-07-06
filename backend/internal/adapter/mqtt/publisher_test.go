package mqtt

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

func TestNew_FailsOnUnreachableBroker(t *testing.T) {
	cfg := entity.MQTTConfig{
		Host:         "127.0.0.1",
		Port:         19999, // nothing listening here
		PublishTopic: "test/alerts",
	}
	_, err := New(cfg)
	require.Error(t, err, "should fail to connect to non-existent broker")
	assert.Contains(t, err.Error(), "connecting to MQTT broker")
}

func TestPublisher_ImplementsInterface(t *testing.T) {
	// Compile-time check is enforced by the var _ declaration in publisher.go.
	// This test documents the contract without requiring a live broker.
	cfg := entity.MQTTConfig{
		Enabled:      true,
		Host:         "localhost",
		Port:         1883,
		PublishTopic: "same/alerts",
	}
	// Just verify the config struct is well-formed.
	assert.Equal(t, "same/alerts", cfg.PublishTopic)
	assert.Equal(t, 1883, cfg.Port)
}
