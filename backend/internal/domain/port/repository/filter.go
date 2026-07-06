package repository

import (
	"context"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

type FilterRepository interface {
	Get(ctx context.Context) (entity.AlertFilter, error)
	Update(ctx context.Context, filter entity.AlertFilter) error
}

type SDRConfigRepository interface {
	Get(ctx context.Context) (entity.SDRDeviceConfig, error)
	Update(ctx context.Context, cfg entity.SDRDeviceConfig) error
}
