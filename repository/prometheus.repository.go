package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type PrometheusRepository interface {
	Upsert(ctx context.Context, data entity.PrometheusUpsert) error
}

type prometheusRepository struct {
	db *sqlx.DB
}

func NewPrometheusRepository(db *sqlx.DB) *prometheusRepository {
	return &prometheusRepository{db}
}

func (r *prometheusRepository) Upsert(ctx context.Context, data entity.PrometheusUpsert) error {
	query := `INSERT INTO prometheus_devices (region, state, ip, idx, shell, card, port, status)
	VALUES (:region, :state, :ip, :idx, :shell, :card, :port, :status)
	ON CONFLICT (ip, idx, shell, card, port) DO UPDATE SET
		region = EXCLUDED.region,
		state = EXCLUDED.state,
		idx = EXCLUDED.idx,
		shell = EXCLUDED.shell,
		card = EXCLUDED.card,
		port = EXCLUDED.port,
		status =  EXCLUDED.status,
		created_at = NOW()`
	_, err := r.db.NamedExecContext(ctx, query, data)
	if err != nil {
		return fmt.Errorf("failed to upsert prometheus device: %w", err)
	}
	return nil
}
