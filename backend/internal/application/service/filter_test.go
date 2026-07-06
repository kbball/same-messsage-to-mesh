package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

// stubSDRConfigRepo implements portrepo.SDRConfigRepository.
type stubSDRConfigRepo struct {
	cfg entity.SDRDeviceConfig
	err error
}

func (s *stubSDRConfigRepo) Get(_ context.Context) (entity.SDRDeviceConfig, error) {
	return s.cfg, s.err
}
func (s *stubSDRConfigRepo) Update(_ context.Context, cfg entity.SDRDeviceConfig) error {
	s.cfg = cfg
	return s.err
}

// stubMQTTConfigRepo implements portrepo.MQTTConfigRepository.
type stubMQTTConfigRepo struct {
	cfg entity.MQTTConfig
	err error
}

func (s *stubMQTTConfigRepo) Get(_ context.Context) (entity.MQTTConfig, error) {
	return s.cfg, s.err
}
func (s *stubMQTTConfigRepo) Update(_ context.Context, cfg entity.MQTTConfig) error {
	s.cfg = cfg
	return s.err
}

func newFilterSvc(filter entity.AlertFilter, sdr entity.SDRDeviceConfig, mqtt entity.MQTTConfig) *FilterService {
	return NewFilterService(
		&stubFilterRepo{filter: filter},
		&stubSDRConfigRepo{cfg: sdr},
		&stubMQTTConfigRepo{cfg: mqtt},
	)
}

func TestFilterService_GetFilter(t *testing.T) {
	want := entity.AlertFilter{StateCodes: []string{"13"}}
	svc := newFilterSvc(want, entity.SDRDeviceConfig{}, entity.MQTTConfig{})
	got, err := svc.GetFilter(context.Background())
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestFilterService_UpdateFilter(t *testing.T) {
	repo := &stubFilterRepo{}
	svc := NewFilterService(repo, &stubSDRConfigRepo{}, &stubMQTTConfigRepo{})
	err := svc.UpdateFilter(context.Background(), entity.AlertFilter{EventCodes: []string{"TOR"}})
	require.NoError(t, err)
	assert.Equal(t, []string{"TOR"}, repo.filter.EventCodes)
}

func TestFilterService_UpdateFilter_Error(t *testing.T) {
	repo := &errFilterRepo{err: errors.New("db error")}
	svc := NewFilterService(repo, &stubSDRConfigRepo{}, &stubMQTTConfigRepo{})
	err := svc.UpdateFilter(context.Background(), entity.AlertFilter{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "updating filter")
}

// errFilterRepo is a FilterRepository stub that always errors on Update.
type errFilterRepo struct {
	stubFilterRepo
	err error
}

func (e *errFilterRepo) Update(_ context.Context, _ entity.AlertFilter) error { return e.err }

func TestFilterService_GetSDRConfig(t *testing.T) {
	want := entity.SDRDeviceConfig{DevicePath: "/dev/rtl0", Frequency: 162550000}
	svc := newFilterSvc(entity.AlertFilter{}, want, entity.MQTTConfig{})
	got, err := svc.GetSDRConfig(context.Background())
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestFilterService_UpdateSDRConfig(t *testing.T) {
	sdrRepo := &stubSDRConfigRepo{}
	svc := NewFilterService(&stubFilterRepo{}, sdrRepo, &stubMQTTConfigRepo{})
	cfg := entity.SDRDeviceConfig{DevicePath: "/dev/rtl0", Frequency: 162400000}
	err := svc.UpdateSDRConfig(context.Background(), cfg)
	require.NoError(t, err)
	assert.Equal(t, cfg, sdrRepo.cfg)
}

func TestFilterService_UpdateSDRConfig_Error(t *testing.T) {
	sdrRepo := &stubSDRConfigRepo{err: errors.New("db error")}
	svc := NewFilterService(&stubFilterRepo{}, sdrRepo, &stubMQTTConfigRepo{})
	err := svc.UpdateSDRConfig(context.Background(), entity.SDRDeviceConfig{DevicePath: "/dev/x"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "updating SDR config")
}

func TestFilterService_GetMQTTConfig(t *testing.T) {
	want := entity.MQTTConfig{Enabled: true, Host: "broker.local", Port: 1883, PublishTopic: "same/alerts"}
	svc := newFilterSvc(entity.AlertFilter{}, entity.SDRDeviceConfig{}, want)
	got, err := svc.GetMQTTConfig(context.Background())
	require.NoError(t, err)
	assert.Equal(t, want, got)
}

func TestFilterService_UpdateMQTTConfig(t *testing.T) {
	mqttRepo := &stubMQTTConfigRepo{}
	svc := NewFilterService(&stubFilterRepo{}, &stubSDRConfigRepo{}, mqttRepo)
	cfg := entity.MQTTConfig{Enabled: true, Host: "broker.local", Port: 1883, PublishTopic: "alerts"}
	err := svc.UpdateMQTTConfig(context.Background(), cfg)
	require.NoError(t, err)
	assert.Equal(t, cfg, mqttRepo.cfg)
}

func TestFilterService_UpdateMQTTConfig_Error(t *testing.T) {
	mqttRepo := &stubMQTTConfigRepo{err: errors.New("db error")}
	svc := NewFilterService(&stubFilterRepo{}, &stubSDRConfigRepo{}, mqttRepo)
	err := svc.UpdateMQTTConfig(context.Background(), entity.MQTTConfig{Host: "x"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "updating MQTT config")
}
