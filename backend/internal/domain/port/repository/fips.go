package repository

import (
	"context"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
)

type FIPSRepository interface {
	ListStates(ctx context.Context) ([]entity.FIPSCode, error)
	ListByState(ctx context.Context, stateCode string) ([]entity.FIPSCode, error)
	GetByFIPS(ctx context.Context, fips string) (entity.FIPSCode, error)
	Upsert(ctx context.Context, codes []entity.FIPSCode) error
	Count(ctx context.Context) (int, error)
}
