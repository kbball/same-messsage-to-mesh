package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

func (h *Handler) getFilter(w http.ResponseWriter, r *http.Request) {
	filter, err := h.filters.GetFilter(r.Context())
	if err != nil {
		slog.Error("failed to get filter config", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get filter config")
		return
	}
	writeJSON(w, http.StatusOK, filter)
}

func (h *Handler) updateFilter(w http.ResponseWriter, r *http.Request) {
	var filter entity.AlertFilter
	if err := json.NewDecoder(r.Body).Decode(&filter); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.filters.UpdateFilter(r.Context(), filter); err != nil {
		slog.Error("failed to update filter config", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update filter config")
		return
	}
	updated, err := h.filters.GetFilter(r.Context())
	if err != nil {
		slog.Error("failed to read updated filter", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to read updated filter")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) getSDRConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := h.filters.GetSDRConfig(r.Context())
	if err != nil {
		slog.Error("failed to get SDR config", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to get SDR config")
		return
	}
	writeJSON(w, http.StatusOK, cfg)
}

func (h *Handler) updateSDRConfig(w http.ResponseWriter, r *http.Request) {
	var cfg entity.SDRDeviceConfig
	if err := json.NewDecoder(r.Body).Decode(&cfg); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if cfg.DevicePath == "" {
		writeError(w, http.StatusBadRequest, "device_path is required")
		return
	}
	if cfg.Frequency <= 0 {
		writeError(w, http.StatusBadRequest, "frequency must be positive")
		return
	}
	if err := h.filters.UpdateSDRConfig(r.Context(), cfg); err != nil {
		slog.Error("failed to update SDR config", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to update SDR config")
		return
	}
	updated, err := h.filters.GetSDRConfig(r.Context())
	if err != nil {
		slog.Error("failed to read updated SDR config", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to read updated SDR config")
		return
	}
	if h.restartSDR != nil {
		if err := h.restartSDR(updated); err != nil {
			slog.Warn("SDR pipeline restart failed after config update", "error", err)
			writeJSON(w, http.StatusOK, map[string]string{
				"warning": "config saved but SDR pipeline restart failed",
			})
			return
		}
	}
	writeJSON(w, http.StatusOK, updated)
}
