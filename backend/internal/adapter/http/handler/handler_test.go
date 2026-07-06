package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kbball/same-message-to-mesh/backend/internal/application/service"
	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

// ---- stub repos -------------------------------------------------------

type memFilterRepo struct{ f entity.AlertFilter }

func (r *memFilterRepo) Get(_ context.Context) (entity.AlertFilter, error) { return r.f, nil }
func (r *memFilterRepo) Update(_ context.Context, f entity.AlertFilter) error {
	r.f = f
	return nil
}

type memSDRRepo struct{ cfg entity.SDRDeviceConfig }

func (r *memSDRRepo) Get(_ context.Context) (entity.SDRDeviceConfig, error) { return r.cfg, nil }
func (r *memSDRRepo) Update(_ context.Context, c entity.SDRDeviceConfig) error {
	r.cfg = c
	return nil
}

type memMQTTRepo struct{ cfg entity.MQTTConfig }

func (r *memMQTTRepo) Get(_ context.Context) (entity.MQTTConfig, error) { return r.cfg, nil }
func (r *memMQTTRepo) Update(_ context.Context, c entity.MQTTConfig) error {
	r.cfg = c
	return nil
}

type memAlertRepo struct{ alerts []entity.SAMEAlert }

func (r *memAlertRepo) Create(_ context.Context, a entity.SAMEAlert) (entity.SAMEAlert, error) {
	a.ID = int64(len(r.alerts) + 1)
	a.ReceivedAt = time.Now()
	r.alerts = append(r.alerts, a)
	return a, nil
}
func (r *memAlertRepo) List(_ context.Context, limit int) ([]entity.SAMEAlert, error) {
	if limit > len(r.alerts) {
		return r.alerts, nil
	}
	return r.alerts[:limit], nil
}
func (r *memAlertRepo) MarkPublished(_ context.Context, _ int64) error { return nil }

type memFIPSRepo struct{ states []entity.FIPSCode }

func (r *memFIPSRepo) ListStates(_ context.Context) ([]entity.FIPSCode, error) {
	return r.states, nil
}
func (r *memFIPSRepo) ListByState(_ context.Context, code string) ([]entity.FIPSCode, error) {
	var out []entity.FIPSCode
	for _, f := range r.states {
		if f.StateCode == code {
			out = append(out, f)
		}
	}
	return out, nil
}
func (r *memFIPSRepo) GetByFIPS(_ context.Context, fips string) (entity.FIPSCode, error) {
	return entity.FIPSCode{StateCode: fips[:2], CountyCode: fips[2:]}, nil
}
func (r *memFIPSRepo) Upsert(_ context.Context, codes []entity.FIPSCode) error {
	r.states = append(r.states, codes...)
	return nil
}
func (r *memFIPSRepo) Count(_ context.Context) (int, error) { return len(r.states), nil }

type memECRepo struct{ codes []entity.EventCode }

func (r *memECRepo) List(_ context.Context) ([]entity.EventCode, error) { return r.codes, nil }
func (r *memECRepo) Get(_ context.Context, code string) (entity.EventCode, error) {
	for _, c := range r.codes {
		if c.Code == code {
			return c, nil
		}
	}
	return entity.EventCode{Code: code}, nil
}
func (r *memECRepo) Upsert(_ context.Context, codes []entity.EventCode) error {
	r.codes = append(r.codes, codes...)
	return nil
}

type stubFetcher struct {
	fips   []entity.FIPSCode
	events []entity.EventCode
	err    error
}

func (f *stubFetcher) FetchFIPS(_ context.Context) ([]entity.FIPSCode, error) {
	return f.fips, f.err
}
func (f *stubFetcher) FetchEventCodes(_ context.Context) ([]entity.EventCode, error) {
	return f.events, f.err
}

type stubPublisher struct{ err error }

func (p *stubPublisher) Publish(_ context.Context, _ entity.SAMEAlert, _ string) error {
	return p.err
}

// nullSSE satisfies sse.Publisher without doing anything.
type nullSSE struct{}

func (n *nullSSE) Publish(_ string, _ any) {}

// ---- error stub repos -------------------------------------------------

type failAlertRepo struct{ memAlertRepo }

func (r *failAlertRepo) List(_ context.Context, _ int) ([]entity.SAMEAlert, error) {
	return nil, errors.New("db error")
}

type failFilterRepo struct{ memFilterRepo }

func (r *failFilterRepo) Get(_ context.Context) (entity.AlertFilter, error) {
	return entity.AlertFilter{}, errors.New("db error")
}

type failSDRRepo struct{ memSDRRepo }

func (r *failSDRRepo) Get(_ context.Context) (entity.SDRDeviceConfig, error) {
	return entity.SDRDeviceConfig{}, errors.New("db error")
}

type failMQTTRepo struct{ memMQTTRepo }

func (r *failMQTTRepo) Get(_ context.Context) (entity.MQTTConfig, error) {
	return entity.MQTTConfig{}, errors.New("db error")
}

type failFIPSRepo struct{ memFIPSRepo }

func (r *failFIPSRepo) ListStates(_ context.Context) ([]entity.FIPSCode, error) {
	return nil, errors.New("db error")
}
func (r *failFIPSRepo) ListByState(_ context.Context, _ string) ([]entity.FIPSCode, error) {
	return nil, errors.New("db error")
}
func (r *failFIPSRepo) Count(_ context.Context) (int, error) {
	return 0, errors.New("db error")
}

type failECRepo struct{ memECRepo }

func (r *failECRepo) List(_ context.Context) ([]entity.EventCode, error) {
	return nil, errors.New("db error")
}

type failFetcher struct{}

func (f *failFetcher) FetchFIPS(_ context.Context) ([]entity.FIPSCode, error) {
	return nil, errors.New("network error")
}
func (f *failFetcher) FetchEventCodes(_ context.Context) ([]entity.EventCode, error) {
	return nil, errors.New("network error")
}

// updateFailFilterRepo: Get succeeds, Update fails.
type updateFailFilterRepo struct{ memFilterRepo }

func (r *updateFailFilterRepo) Update(_ context.Context, _ entity.AlertFilter) error {
	return errors.New("db write error")
}

// updateFailSDRRepo: Get succeeds, Update fails.
type updateFailSDRRepo struct{ memSDRRepo }

func (r *updateFailSDRRepo) Update(_ context.Context, _ entity.SDRDeviceConfig) error {
	return errors.New("db write error")
}

// updateFailMQTTRepo: Get succeeds, Update fails.
type updateFailMQTTRepo struct{ memMQTTRepo }

func (r *updateFailMQTTRepo) Update(_ context.Context, _ entity.MQTTConfig) error {
	return errors.New("db write error")
}

// ---- test helpers -----------------------------------------------------

func newTestHandler() (*Handler, *memAlertRepo, *memFilterRepo, *memFIPSRepo, *memECRepo) {
	alertRepo := &memAlertRepo{}
	filterRepo := &memFilterRepo{}
	sdrRepo := &memSDRRepo{cfg: entity.SDRDeviceConfig{DevicePath: "/dev/rtl0", Frequency: 162550000}}
	mqttRepo := &memMQTTRepo{cfg: entity.MQTTConfig{Host: "localhost", Port: 1883}}
	fipsRepo := &memFIPSRepo{states: []entity.FIPSCode{
		{StateCode: "13", CountyCode: "121", StateName: "Georgia", CountyName: "Fulton County"},
	}}
	ecRepo := &memECRepo{codes: []entity.EventCode{
		{Code: "TOR", Description: "Tornado Warning"},
	}}
	fetcher := &stubFetcher{
		fips:   []entity.FIPSCode{{StateCode: "13", CountyCode: "121"}},
		events: []entity.EventCode{{Code: "TOR"}},
	}

	alertSvc := service.NewAlertService(alertRepo, filterRepo, fipsRepo, ecRepo, nil)
	filterSvc := service.NewFilterService(filterRepo, sdrRepo, mqttRepo)
	refSvc := service.NewReferenceDataService(fipsRepo, ecRepo, fetcher)

	h := New(alertSvc, filterSvc, refSvc, &nullSSE{})
	return h, alertRepo, filterRepo, fipsRepo, ecRepo
}

func do(t *testing.T, h *Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var buf bytes.Buffer
	if body != nil {
		require.NoError(t, json.NewEncoder(&buf).Encode(body))
	}
	req := httptest.NewRequest(method, path, &buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	mux := http.NewServeMux()
	h.Register(mux)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	return rr
}

// ---- alerts -----------------------------------------------------------

func TestHandler_ListAlerts_Empty(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodGet, "/api/alerts", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got []entity.SAMEAlert
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Empty(t, got)
}

func TestHandler_ListAlerts_WithData(t *testing.T) {
	h, alertRepo, _, _, _ := newTestHandler()
	alertRepo.alerts = []entity.SAMEAlert{
		{ID: 1, EventCode: "TOR", ReceivedAt: time.Now()},
	}
	rr := do(t, h, http.MethodGet, "/api/alerts", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got []entity.SAMEAlert
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Len(t, got, 1)
	assert.Equal(t, "TOR", got[0].EventCode)
}

// ---- filters ----------------------------------------------------------

func TestHandler_GetFilter(t *testing.T) {
	h, _, filterRepo, _, _ := newTestHandler()
	filterRepo.f = entity.AlertFilter{StateCodes: []string{"13"}}
	rr := do(t, h, http.MethodGet, "/api/filters", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got entity.AlertFilter
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Equal(t, []string{"13"}, got.StateCodes)
}

func TestHandler_UpdateFilter(t *testing.T) {
	h, _, filterRepo, _, _ := newTestHandler()
	rr := do(t, h, http.MethodPut, "/api/filters", entity.AlertFilter{EventCodes: []string{"TOR"}})
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, []string{"TOR"}, filterRepo.f.EventCodes)
}

func TestHandler_UpdateFilter_BadBody(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	req := httptest.NewRequest(http.MethodPut, "/api/filters", bytes.NewBufferString("not-json"))
	mux := http.NewServeMux()
	h.Register(mux)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ---- SDR config -------------------------------------------------------

func TestHandler_GetSDRConfig(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodGet, "/api/sdr-config", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got entity.SDRDeviceConfig
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Equal(t, "/dev/rtl0", got.DevicePath)
}

func TestHandler_UpdateSDRConfig(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodPut, "/api/sdr-config", entity.SDRDeviceConfig{DevicePath: "/dev/rtl1", Frequency: 162400000})
	assert.Equal(t, http.StatusOK, rr.Code)
	var got entity.SDRDeviceConfig
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Equal(t, "/dev/rtl1", got.DevicePath)
}

func TestHandler_UpdateSDRConfig_MissingDevicePath(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodPut, "/api/sdr-config", entity.SDRDeviceConfig{Frequency: 162550000})
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_UpdateSDRConfig_InvalidFrequency(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodPut, "/api/sdr-config", entity.SDRDeviceConfig{DevicePath: "/dev/rtl0"})
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// ---- MQTT config ------------------------------------------------------

func TestHandler_GetMQTTConfig(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodGet, "/api/mqtt-config", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got entity.MQTTConfig
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Equal(t, "localhost", got.Host)
}

func TestHandler_UpdateMQTTConfig(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodPut, "/api/mqtt-config", entity.MQTTConfig{Host: "broker.local", Port: 1883, PublishTopic: "same/alerts"})
	assert.Equal(t, http.StatusOK, rr.Code)
	var got entity.MQTTConfig
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Equal(t, "broker.local", got.Host)
}

func TestHandler_UpdateMQTTConfig_MissingHost(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodPut, "/api/mqtt-config", entity.MQTTConfig{Port: 1883})
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestHandler_UpdateMQTTConfig_ReconnectError(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	h.reconnectMQTT = func(_ entity.MQTTConfig) error { return errors.New("connection refused") }
	rr := do(t, h, http.MethodPut, "/api/mqtt-config", entity.MQTTConfig{Host: "broker.local", Port: 1883})
	assert.Equal(t, http.StatusOK, rr.Code)
	var got map[string]string
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Contains(t, got["warning"], "MQTT reconnect failed")
}

func TestHandler_TestMQTTPublish_NoPublisher(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodPost, "/api/mqtt-config/test", nil)
	assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
}

func TestHandler_TestMQTTPublish_Success(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	h.mqttPublisher = &stubPublisher{}
	rr := do(t, h, http.MethodPost, "/api/mqtt-config/test", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got map[string]bool
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.True(t, got["ok"])
}

func TestHandler_TestMQTTPublish_PublishError(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	h.mqttPublisher = &stubPublisher{err: errors.New("broker timeout")}
	rr := do(t, h, http.MethodPost, "/api/mqtt-config/test", nil)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// ---- reference data ---------------------------------------------------

func TestHandler_ListStates(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodGet, "/api/reference/states", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got []entity.FIPSCode
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.NotEmpty(t, got)
}

func TestHandler_ListCounties(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodGet, "/api/reference/counties/13", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got []entity.FIPSCode
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.NotEmpty(t, got)
}

func TestHandler_ListEventCodes(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodGet, "/api/reference/event-codes", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got []entity.EventCode
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.NotEmpty(t, got)
}

func TestHandler_FIPSCount(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodGet, "/api/reference/fips/count", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got map[string]int
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.GreaterOrEqual(t, got["count"], 0)
}

func TestHandler_RefreshFIPS(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodPost, "/api/reference/fips/refresh", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got map[string]int
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Equal(t, 1, got["updated"])
}

func TestHandler_RefreshEventCodes(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	rr := do(t, h, http.MethodPost, "/api/reference/event-codes/refresh", nil)
	assert.Equal(t, http.StatusOK, rr.Code)
	var got map[string]int
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Equal(t, 1, got["updated"])
}

// ---- error paths -------------------------------------------------------

func newErrHandler() *Handler {
	alertSvc := service.NewAlertService(&failAlertRepo{}, &failFilterRepo{}, &failFIPSRepo{}, &failECRepo{}, nil)
	filterSvc := service.NewFilterService(&failFilterRepo{}, &failSDRRepo{}, &failMQTTRepo{})
	refSvc := service.NewReferenceDataService(&failFIPSRepo{}, &failECRepo{}, &failFetcher{})
	return New(alertSvc, filterSvc, refSvc, &nullSSE{})
}

func TestHandler_ListAlerts_Error(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodGet, "/api/alerts", nil)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_GetFilter_Error(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodGet, "/api/filters", nil)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// UpdateFilter with failRepo: Update succeeds but re-read fails.
func TestHandler_UpdateFilter_RereadError(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodPut, "/api/filters", entity.AlertFilter{EventCodes: []string{"TOR"}})
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_UpdateFilter_UpdateError(t *testing.T) {
	filterSvc := service.NewFilterService(&updateFailFilterRepo{}, &memSDRRepo{}, &memMQTTRepo{})
	refSvc := service.NewReferenceDataService(&memFIPSRepo{}, &memECRepo{}, &stubFetcher{})
	alertSvc := service.NewAlertService(&memAlertRepo{}, &memFilterRepo{}, &memFIPSRepo{}, &memECRepo{}, nil)
	h := New(alertSvc, filterSvc, refSvc, &nullSSE{})
	rr := do(t, h, http.MethodPut, "/api/filters", entity.AlertFilter{EventCodes: []string{"TOR"}})
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_GetSDRConfig_Error(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodGet, "/api/sdr-config", nil)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// UpdateSDRConfig with failRepo: Update succeeds but re-read fails.
func TestHandler_UpdateSDRConfig_RereadError(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodPut, "/api/sdr-config", entity.SDRDeviceConfig{DevicePath: "/dev/rtl0", Frequency: 162550000})
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_UpdateSDRConfig_UpdateError(t *testing.T) {
	filterSvc := service.NewFilterService(&memFilterRepo{}, &updateFailSDRRepo{}, &memMQTTRepo{})
	refSvc := service.NewReferenceDataService(&memFIPSRepo{}, &memECRepo{}, &stubFetcher{})
	alertSvc := service.NewAlertService(&memAlertRepo{}, &memFilterRepo{}, &memFIPSRepo{}, &memECRepo{}, nil)
	h := New(alertSvc, filterSvc, refSvc, &nullSSE{})
	rr := do(t, h, http.MethodPut, "/api/sdr-config", entity.SDRDeviceConfig{DevicePath: "/dev/rtl0", Frequency: 162550000})
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_UpdateSDRConfig_RestartCalled(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	var restarted entity.SDRDeviceConfig
	h = h.WithSDR(func(cfg entity.SDRDeviceConfig) error {
		restarted = cfg
		return nil
	})
	rr := do(t, h, http.MethodPut, "/api/sdr-config", entity.SDRDeviceConfig{DevicePath: "/dev/rtl0", Frequency: 162550000})
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "/dev/rtl0", restarted.DevicePath)
	assert.Equal(t, int64(162550000), restarted.Frequency)
}

func TestHandler_UpdateSDRConfig_RestartError(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	h = h.WithSDR(func(_ entity.SDRDeviceConfig) error {
		return errors.New("rtl_fm not found")
	})
	rr := do(t, h, http.MethodPut, "/api/sdr-config", entity.SDRDeviceConfig{DevicePath: "/dev/rtl0", Frequency: 162550000})
	assert.Equal(t, http.StatusOK, rr.Code)
	var got map[string]string
	require.NoError(t, json.NewDecoder(rr.Body).Decode(&got))
	assert.Contains(t, got["warning"], "SDR pipeline restart failed")
}

func TestHandler_GetMQTTConfig_Error(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodGet, "/api/mqtt-config", nil)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// UpdateMQTTConfig with failRepo: Update succeeds but re-read fails.
func TestHandler_UpdateMQTTConfig_RereadError(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodPut, "/api/mqtt-config", entity.MQTTConfig{Host: "broker.local", Port: 1883})
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_UpdateMQTTConfig_UpdateError(t *testing.T) {
	filterSvc := service.NewFilterService(&memFilterRepo{}, &memSDRRepo{}, &updateFailMQTTRepo{})
	refSvc := service.NewReferenceDataService(&memFIPSRepo{}, &memECRepo{}, &stubFetcher{})
	alertSvc := service.NewAlertService(&memAlertRepo{}, &memFilterRepo{}, &memFIPSRepo{}, &memECRepo{}, nil)
	h := New(alertSvc, filterSvc, refSvc, &nullSSE{})
	rr := do(t, h, http.MethodPut, "/api/mqtt-config", entity.MQTTConfig{Host: "broker.local", Port: 1883})
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_ListStates_Error(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodGet, "/api/reference/states", nil)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_ListCounties_Error(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodGet, "/api/reference/counties/13", nil)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_ListEventCodes_Error(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodGet, "/api/reference/event-codes", nil)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_FIPSCount_Error(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodGet, "/api/reference/fips/count", nil)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_RefreshFIPS_Error(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodPost, "/api/reference/fips/refresh", nil)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestHandler_RefreshEventCodes_Error(t *testing.T) {
	h := newErrHandler()
	rr := do(t, h, http.MethodPost, "/api/reference/event-codes/refresh", nil)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// ---- middleware / wiring -----------------------------------------------

func TestLoggingMiddleware(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
	mw := LoggingMiddleware(inner)
	rr := httptest.NewRecorder()
	mw.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/api/alerts", nil))
	assert.Equal(t, http.StatusTeapot, rr.Code)
}

func TestHandler_WithMQTT(t *testing.T) {
	h, _, _, _, _ := newTestHandler()
	pub := &stubPublisher{}
	reconnect := func(_ entity.MQTTConfig) error { return nil }
	h2 := h.WithMQTT(pub, reconnect)
	assert.Equal(t, h, h2) // returns same pointer
	assert.NotNil(t, h.mqttPublisher)
	assert.NotNil(t, h.reconnectMQTT)
}
