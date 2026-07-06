package handler

import (
	"net/http"

	"github.com/kbball/same-message-to-mesh/backend/internal/adapter/sse"
	"github.com/kbball/same-message-to-mesh/backend/internal/application/service"
	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

// Handler holds all application services and registers HTTP routes.
type Handler struct {
	alerts        *service.AlertService
	filters       *service.FilterService
	refData       *service.ReferenceDataService
	stream        sse.Publisher
	mqttPublisher publisherGetter
	reconnectMQTT func(entity.MQTTConfig) error
	restartSDR    func(entity.SDRDeviceConfig) error
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

// WithMQTT wires a live publisher and a reconnect callback into the handler.
func (h *Handler) WithMQTT(pub publisherGetter, reconnect func(entity.MQTTConfig) error) *Handler {
	h.mqttPublisher = pub
	h.reconnectMQTT = reconnect
	return h
}

// WithSDR wires a restart callback that is called after SDR config is saved.
func (h *Handler) WithSDR(restart func(entity.SDRDeviceConfig) error) *Handler {
	h.restartSDR = restart
	return h
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

	// MQTT config
	mux.HandleFunc("GET /api/mqtt-config", h.getMQTTConfig)
	mux.HandleFunc("PUT /api/mqtt-config", h.updateMQTTConfig)
	mux.HandleFunc("POST /api/mqtt-config/test", h.testMQTTPublish)

	// Reference data
	mux.HandleFunc("GET /api/reference/states", h.listStates)
	mux.HandleFunc("GET /api/reference/counties/{stateCode}", h.listCounties)
	mux.HandleFunc("GET /api/reference/event-codes", h.listEventCodes)
	mux.HandleFunc("POST /api/reference/fips/refresh", h.refreshFIPS)
	mux.HandleFunc("POST /api/reference/event-codes/refresh", h.refreshEventCodes)
	mux.HandleFunc("GET /api/reference/fips/count", h.fipsCount)
}
