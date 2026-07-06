package repository

import (
	"context"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

type AlertRepository interface {
	Create(ctx context.Context, alert entity.SAMEAlert) (entity.SAMEAlert, error)
	List(ctx context.Context, limit int) ([]entity.SAMEAlert, error)
	MarkPublished(ctx context.Context, id int64) error
}
