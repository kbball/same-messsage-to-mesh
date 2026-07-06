package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
	portrepo "github.com/kbball/same-message-to-mesh/backend/internal/domain/port/repository"
)

type FilterRepo struct {
	db *sql.DB
}

func NewFilterRepo(db *sql.DB) *FilterRepo { return &FilterRepo{db: db} }

var _ portrepo.FilterRepository = (*FilterRepo)(nil)

func (r *FilterRepo) Get(ctx context.Context) (entity.AlertFilter, error) {
	var f entity.AlertFilter
	err := r.db.QueryRowContext(ctx,
		`SELECT state_codes, fips_codes, event_codes, updated_at FROM filter_config WHERE id = 1`).
		Scan(pq.Array(&f.StateCodes), pq.Array(&f.FIPSCodes), pq.Array(&f.EventCodes), &f.UpdatedAt)
	if err != nil {
		return entity.AlertFilter{}, fmt.Errorf("getting filter config: %w", err)
	}
	return f, nil
}

func (r *FilterRepo) Update(ctx context.Context, filter entity.AlertFilter) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE filter_config SET state_codes = $1, fips_codes = $2, event_codes = $3, updated_at = NOW() WHERE id = 1`,
		pq.Array(filter.StateCodes), pq.Array(filter.FIPSCodes), pq.Array(filter.EventCodes))
	if err != nil {
		return fmt.Errorf("updating filter config: %w", err)
	}
	return nil
}

type SDRConfigRepo struct {
	db *sql.DB
}

func NewSDRConfigRepo(db *sql.DB) *SDRConfigRepo { return &SDRConfigRepo{db: db} }

var _ portrepo.SDRConfigRepository = (*SDRConfigRepo)(nil)

func (r *SDRConfigRepo) Get(ctx context.Context) (entity.SDRDeviceConfig, error) {
	var c entity.SDRDeviceConfig
	err := r.db.QueryRowContext(ctx,
		`SELECT device_path, frequency, updated_at FROM sdr_config WHERE id = 1`).
		Scan(&c.DevicePath, &c.Frequency, &c.UpdatedAt)
	if err != nil {
		return entity.SDRDeviceConfig{}, fmt.Errorf("getting SDR config: %w", err)
	}
	return c, nil
}

func (r *SDRConfigRepo) Update(ctx context.Context, cfg entity.SDRDeviceConfig) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sdr_config SET device_path = $1, frequency = $2, updated_at = NOW() WHERE id = 1`,
		cfg.DevicePath, cfg.Frequency)
	if err != nil {
		return fmt.Errorf("updating SDR config: %w", err)
	}
	return nil
}
