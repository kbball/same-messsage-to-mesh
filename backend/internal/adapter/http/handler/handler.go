package handler

import (
	"net/http"

	"github.com/kbball/same-message-to-mesh/backend/internal/adapter/sse"
	"github.com/kbball/same-message-to-mesh/backend/internal/application/service"
)

// Handler holds all application services and registers HTTP routes.
type Handler struct {
	alerts    *service.AlertService
	filters   *service.FilterService
	refData   *service.ReferenceDataService
	stream    sse.Publisher
}

func New(
	alerts *service.AlertService,
	filters *service.FilterService,
	refData *service.ReferenceDataService,
	stream sse.Publisher,
) *Handler {
	return &Handler{
		alerts:  alerts,
		filters: filters,
		refData: refData,
		stream:  stream,
	}
}

// Register wires all API routes onto mux.
func (h *Handler) Register(mux *http.ServeMux) {
	// Alerts
	mux.HandleFunc("GET /api/alerts", h.listAlerts)

	// Filters
	mux.HandleFunc("GET /api/filters", h.getFilter)
	mux.HandleFunc("PUT /api/filters", h.updateFilter)

	// SDR config
	mux.HandleFunc("GET /api/sdr-config", h.getSDRConfig)
	mux.HandleFunc("PUT /api/sdr-config", h.updateSDRConfig)

	// Reference data
	mux.HandleFunc("GET /api/reference/states", h.listStates)
	mux.HandleFunc("GET /api/reference/counties/{stateCode}", h.listCounties)
	mux.HandleFunc("GET /api/reference/event-codes", h.listEventCodes)
	mux.HandleFunc("POST /api/reference/fips/refresh", h.refreshFIPS)
	mux.HandleFunc("POST /api/reference/event-codes/refresh", h.refreshEventCodes)
	mux.HandleFunc("GET /api/reference/fips/count", h.fipsCount)
}
