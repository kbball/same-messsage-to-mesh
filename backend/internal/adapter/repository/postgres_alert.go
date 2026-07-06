package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/lib/pq"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
	portrepo "github.com/kbball/same-message-to-mesh/backend/internal/domain/port/repository"
)

type AlertRepo struct {
	db *sql.DB
}

func NewAlertRepo(db *sql.DB) *AlertRepo { return &AlertRepo{db: db} }

var _ portrepo.AlertRepository = (*AlertRepo)(nil)

func (r *AlertRepo) Create(ctx context.Context, alert entity.SAMEAlert) (entity.SAMEAlert, error) {
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO same_alerts (originator, event_code, fips_codes, purge_time, issue_time, call_sign, raw_message)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 RETURNING id, received_at, published`,
		alert.Originator,
		alert.EventCode,
		pq.Array(alert.FIPSCodes),
		alert.PurgeTime,
		alert.IssueTime,
		alert.CallSign,
		alert.RawMessage,
	).Scan(&alert.ID, &alert.ReceivedAt, &alert.Published)
	if err != nil {
		return entity.SAMEAlert{}, fmt.Errorf("creating alert: %w", err)
	}
	return alert, nil
}

func (r *AlertRepo) List(ctx context.Context, limit int) ([]entity.SAMEAlert, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, received_at, originator, event_code, fips_codes, purge_time, issue_time, call_sign, raw_message, published
		 FROM same_alerts
		 ORDER BY received_at DESC
		 LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("listing alerts: %w", err)
	}
	defer func() { _ = rows.Close() }()

	var alerts []entity.SAMEAlert
	for rows.Next() {
		var a entity.SAMEAlert
		if err := rows.Scan(
			&a.ID, &a.ReceivedAt, &a.Originator, &a.EventCode,
			pq.Array(&a.FIPSCodes),
			&a.PurgeTime, &a.IssueTime, &a.CallSign, &a.RawMessage, &a.Published,
		); err != nil {
			return nil, fmt.Errorf("scanning alert: %w", err)
		}
		alerts = append(alerts, a)
	}
	return alerts, rows.Err()
}

func (r *AlertRepo) MarkPublished(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `UPDATE same_alerts SET published = true WHERE id = $1`, id)
	return err
}
