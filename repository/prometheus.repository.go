package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type PrometheusRepository interface {
	Upsert(ctx context.Context, data entity.PrometheusUpsert) error
	GponPortsStatus(ctx context.Context) (*entity.PrometheusPortStatus, error)
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

func (r *prometheusRepository) GponPortsStatus(ctx context.Context) (*entity.PrometheusPortStatus, error) {
	var res entity.PrometheusPortStatus
	query := `SELECT
		COUNT(DISTINCT(ip)) AS olts,
		COUNT(DISTINCT(ip, card)) AS cards,
		SUM(CASE WHEN status IN (1, 6) THEN 1 ELSE 0 END) AS gpon_actives,
  		SUM(CASE WHEN status NOT IN (1, 6) THEN 1 ELSE 0 END) AS gpon_inactives,
  		COUNT(*) AS total_gpon
	FROM prometheus_devices;`
	err := r.db.GetContext(ctx, &res, query)
	return &res, err
}
