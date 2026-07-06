package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
	portrepo "github.com/kbball/same-message-to-mesh/backend/internal/domain/port/repository"
)

type EventCodeRepo struct {
	db *sql.DB
}

func NewEventCodeRepo(db *sql.DB) *EventCodeRepo { return &EventCodeRepo{db: db} }

var _ portrepo.EventCodeRepository = (*EventCodeRepo)(nil)

func (r *EventCodeRepo) List(ctx context.Context) ([]entity.EventCode, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT code, description, category, is_warning, updated_at FROM event_codes ORDER BY code`)
	if err != nil {
		return nil, fmt.Errorf("listing event codes: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var codes []entity.EventCode
	for rows.Next() {
		var c entity.EventCode
		if err := rows.Scan(&c.Code, &c.Description, &c.Category, &c.IsWarning, &c.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning event code: %w", err)
		}
		codes = append(codes, c)
	}
	return codes, rows.Err()
}

func (r *EventCodeRepo) Get(ctx context.Context, code string) (entity.EventCode, error) {
	var c entity.EventCode
	err := r.db.QueryRowContext(ctx,
		`SELECT code, description, category, is_warning, updated_at FROM event_codes WHERE code = $1`, code).
		Scan(&c.Code, &c.Description, &c.Category, &c.IsWarning, &c.UpdatedAt)
	if err != nil {
		return entity.EventCode{}, mapNotFound(err)
	}
	return c, nil
}

func (r *EventCodeRepo) Upsert(ctx context.Context, codes []entity.EventCode) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO event_codes (code, description, category, is_warning, updated_at)
		 VALUES ($1, $2, $3, $4, NOW())
		 ON CONFLICT (code) DO UPDATE SET
		   description = EXCLUDED.description,
		   category    = EXCLUDED.category,
		   is_warning  = EXCLUDED.is_warning,
		   updated_at  = NOW()`)
	if err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	for _, c := range codes {
		if _, err := stmt.ExecContext(ctx, c.Code, c.Description, c.Category, c.IsWarning); err != nil {
			return fmt.Errorf("upserting event code %s: %w", c.Code, err)
		}
	}
	return tx.Commit()
}
