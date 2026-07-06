package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
	portsvc "github.com/kbball/same-message-to-mesh/backend/internal/domain/port/service"
)

func (h *Handler) getMQTTConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.filters.GetMQTTConfig(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to get MQTT config")
		return
	}
	writeJSON(w, http.StatusOK, cfg)
}

func (h *Handler) updateMQTTConfig(w http.ResponseWriter, r *http.Request) {
	var cfg entity.MQTTConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if cfg.Host == "" {
		writeError(w, http.StatusBadRequest, "host is required")
		return
	}
	if cfg.Port <= 0 {
		writeError(w, http.StatusBadRequest, "port must be positive")
		return
	}
	if err := h.filters.UpdateMQTTConfig(r.Context(), cfg); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update MQTT config")
		return
	}
	// Reconnect the publisher with the new config.
	if h.reconnectMQTT != nil {
		if err := h.reconnectMQTT(cfg); err != nil {
			// Non-fatal: config is saved, but warn the client.
			writeJSON(w, http.StatusOK, map[string]string{
				"warning": "config saved but MQTT reconnect failed: " + err.Error(),
			})
			return
		}
	}
	updated, err := h.filters.GetMQTTConfig(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read updated MQTT config")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) testMQTTPublish(w http.ResponseWriter, r *http.Request) {
	pub := h.mqttPublisher
	if pub == nil {
		writeError(w, http.StatusServiceUnavailable, "MQTT publisher is not enabled")
		return
	}
	msg := "[TEST] SAME → Mesh connectivity check"
	if err := pub.Publish(context.Background(), entity.SAMEAlert{}, msg); err != nil {
		writeError(w, http.StatusInternalServerError, "MQTT test publish failed: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// publisherGetter is a narrow interface so the handler can call Publish without importing the mqtt package.
type publisherGetter interface {
	Publish(ctx context.Context, alert entity.SAMEAlert, message string) error
}

var _ portsvc.AlertPublisher = (publisherGetter)(nil)
