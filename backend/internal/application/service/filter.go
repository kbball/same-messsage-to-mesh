package service

import (
	"context"
	"fmt"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
	portrepo "github.com/kbball/same-message-to-mesh/backend/internal/domain/port/repository"
)

type FilterService struct {
	filterRepo     portrepo.FilterRepository
	sdrConfigRepo  portrepo.SDRConfigRepository
	mqttConfigRepo portrepo.MQTTConfigRepository
}

func NewFilterService(
	filterRepo portrepo.FilterRepository,
	sdrConfigRepo portrepo.SDRConfigRepository,
	mqttConfigRepo portrepo.MQTTConfigRepository,
) *FilterService {
	return &FilterService{
		filterRepo:     filterRepo,
		sdrConfigRepo:  sdrConfigRepo,
		mqttConfigRepo: mqttConfigRepo,
	}
}

func (s *FilterService) GetFilter(ctx context.Context) (entity.AlertFilter, error) {
	return s.filterRepo.Get(ctx)
}

func (s *FilterService) UpdateFilter(ctx context.Context, filter entity.AlertFilter) error {
	if err := s.filterRepo.Update(ctx, filter); err != nil {
		return fmt.Errorf("updating filter: %w", err)
	}
	return nil
}

func (s *FilterService) GetSDRConfig(ctx context.Context) (entity.SDRDeviceConfig, error) {
	return s.sdrConfigRepo.Get(ctx)
}

func (s *FilterService) UpdateSDRConfig(ctx context.Context, cfg entity.SDRDeviceConfig) error {
	if err := s.sdrConfigRepo.Update(ctx, cfg); err != nil {
		return fmt.Errorf("updating SDR config: %w", err)
	}
	return nil
}

func (s *FilterService) GetMQTTConfig(ctx context.Context) (entity.MQTTConfig, error) {
	return s.mqttConfigRepo.Get(ctx)
}

func (s *FilterService) UpdateMQTTConfig(ctx context.Context, cfg entity.MQTTConfig) error {
	if err := s.mqttConfigRepo.Update(ctx, cfg); err != nil {
		return fmt.Errorf("updating MQTT config: %w", err)
	}
	return nil
}
