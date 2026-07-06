package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kbball/same-message-to-mesh/backend/internal/domain/entity"
	portrepo "github.com/kbball/same-message-to-mesh/backend/internal/domain/port/repository"
)

type MQTTConfigRepo struct {
	db *sql.DB
}

func NewMQTTConfigRepo(db *sql.DB) *MQTTConfigRepo { return &MQTTConfigRepo{db: db} }

var _ portrepo.MQTTConfigRepository = (*MQTTConfigRepo)(nil)

func (r *MQTTConfigRepo) Get(ctx context.Context) (entity.MQTTConfig, error) {
	var c entity.MQTTConfig
	err := r.db.QueryRowContext(ctx,
		`SELECT enabled, host, port, publish_topic, updated_at FROM mqtt_config WHERE id = 1`).
		Scan(&c.Enabled, &c.Host, &c.Port, &c.PublishTopic, &c.UpdatedAt)
	if err != nil {
		return entity.MQTTConfig{}, fmt.Errorf("getting MQTT config: %w", err)
	}
	return c, nil
}

func (r *MQTTConfigRepo) Update(ctx context.Context, cfg entity.MQTTConfig) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE mqtt_config SET enabled = $1, host = $2, port = $3, publish_topic = $4, updated_at = NOW() WHERE id = 1`,
		cfg.Enabled, cfg.Host, cfg.Port, cfg.PublishTopic)
	if err != nil {
		return fmt.Errorf("updating MQTT config: %w", err)
	}
	return nil
}
