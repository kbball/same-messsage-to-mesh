package handler

import (
	"net/http"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

func (h *Handler) listStates(w http.ResponseWriter, r *http.Request) {
	states, err := h.refData.ListStates(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list states")
		return
	}
	if states == nil {
		states = []entity.FIPSCode{}
	}
	writeJSON(w, http.StatusOK, states)
}

func (h *Handler) listCounties(w http.ResponseWriter, r *http.Request) {
	stateCode := r.PathValue("stateCode")
	if stateCode == "" {
		writeError(w, http.StatusBadRequest, "stateCode is required")
		return
	}
	counties, err := h.refData.ListFIPSByState(r.Context(), stateCode)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list counties")
		return
	}
	if counties == nil {
		counties = []entity.FIPSCode{}
	}
	writeJSON(w, http.StatusOK, counties)
}

func (h *Handler) listEventCodes(w http.ResponseWriter, r *http.Request) {
	codes, err := h.refData.ListEventCodes(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list event codes")
		return
	}
	if codes == nil {
		codes = []entity.EventCode{}
	}
	writeJSON(w, http.StatusOK, codes)
}

func (h *Handler) refreshFIPS(w http.ResponseWriter, r *http.Request) {
	count, err := h.refData.RefreshFIPS(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to refresh FIPS codes: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"updated": count})
}

func (h *Handler) refreshEventCodes(w http.ResponseWriter, r *http.Request) {
	count, err := h.refData.RefreshEventCodes(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to refresh event codes: "+err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"updated": count})
}

func (h *Handler) fipsCount(w http.ResponseWriter, r *http.Request) {
	n, err := h.refData.FIPSCount(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to count FIPS codes")
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": n})
}
