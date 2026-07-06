package repository

import (
	"context"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

type EventCodeRepository interface {
	List(ctx context.Context) ([]entity.EventCode, error)
	Get(ctx context.Context, code string) (entity.EventCode, error)
	Upsert(ctx context.Context, codes []entity.EventCode) error
}
