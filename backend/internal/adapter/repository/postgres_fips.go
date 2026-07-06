package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
	portrepo "github.com/kbball/same-message-to-mesh/backend/internal/domain/port/repository"
)

type FIPSRepo struct {
	db *sql.DB
}

func NewFIPSRepo(db *sql.DB) *FIPSRepo { return &FIPSRepo{db: db} }

var _ portrepo.FIPSRepository = (*FIPSRepo)(nil)

func (r *FIPSRepo) ListStates(ctx context.Context) ([]entity.FIPSCode, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DISTINCT ON (state_code) state_code, county_code, state_name, county_name, updated_at
		 FROM fips_codes ORDER BY state_code`)
	if err != nil {
		return nil, fmt.Errorf("listing states: %w", err)
	}
	defer func() { _ = rows.Close() }()
	return scanFIPSRows(rows)
}

func (r *FIPSRepo) ListByState(ctx context.Context, stateCode string) ([]entity.FIPSCode, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT state_code, county_code, state_name, county_name, updated_at
		 FROM fips_codes WHERE state_code = $1 ORDER BY county_name`,
		stateCode)
	if err != nil {
		return nil, fmt.Errorf("listing counties for state %s: %w", stateCode, err)
	}
	defer func() { _ = rows.Close() }()
	return scanFIPSRows(rows)
}

func (r *FIPSRepo) GetByFIPS(ctx context.Context, fips string) (entity.FIPSCode, error) {
	if len(fips) != 5 {
		return entity.FIPSCode{}, fmt.Errorf("invalid FIPS code length: %s", fips)
	}
	var f entity.FIPSCode
	err := r.db.QueryRowContext(ctx,
		`SELECT state_code, county_code, state_name, county_name, updated_at
		 FROM fips_codes WHERE state_code = $1 AND county_code = $2`,
		fips[:2], fips[2:]).
		Scan(&f.StateCode, &f.CountyCode, &f.StateName, &f.CountyName, &f.UpdatedAt)
	if err != nil {
		return entity.FIPSCode{}, mapNotFound(err)
	}
	return f, nil
}

func (r *FIPSRepo) Upsert(ctx context.Context, codes []entity.FIPSCode) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO fips_codes (state_code, county_code, state_name, county_name, updated_at)
		 VALUES ($1, $2, $3, $4, NOW())
		 ON CONFLICT (state_code, county_code) DO UPDATE SET
		   state_name  = EXCLUDED.state_name,
		   county_name = EXCLUDED.county_name,
		   updated_at  = NOW()`)
	if err != nil {
		return fmt.Errorf("preparing statement: %w", err)
	}
	defer func() { _ = stmt.Close() }()

	for _, f := range codes {
		if _, err := stmt.ExecContext(ctx, f.StateCode, f.CountyCode, f.StateName, f.CountyName); err != nil {
			return fmt.Errorf("upserting FIPS %s%s: %w", f.StateCode, f.CountyCode, err)
		}
	}
	return tx.Commit()
}

func (r *FIPSRepo) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM fips_codes`).Scan(&n)
	return n, err
}

func scanFIPSRows(rows *sql.Rows) ([]entity.FIPSCode, error) {
	var codes []entity.FIPSCode
	for rows.Next() {
		var f entity.FIPSCode
		if err := rows.Scan(&f.StateCode, &f.CountyCode, &f.StateName, &f.CountyName, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scanning FIPS code: %w", err)
		}
		codes = append(codes, f)
	}
	return codes, rows.Err()
}
