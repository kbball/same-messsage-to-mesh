package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

func (h *Handler) listAlerts(w http.ResponseWriter, r *http.Request) {
	limit := 100
	if s := r.URL.Query().Get("limit"); s != "" {
		if n, err := strconv.Atoi(s); err == nil && n > 0 {
			limit = n
		}
	}

	alerts, err := h.alerts.List(r.Context(), limit)
	if err != nil {
		slog.Error("failed to list alerts", "error", err)
		writeError(w, http.StatusInternalServerError, "failed to list alerts")
		return
	}
	if alerts == nil {
		alerts = []entity.SAMEAlert{}
	}
	writeJSON(w, http.StatusOK, alerts)
}
