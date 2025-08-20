package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type TrafficRepository interface {
	GetSnmpIndexByODN(ctx context.Context, ip, odn string) ([]string, error)
	GetSnmpIndexByFAT(ctx context.Context, ip, odn, fat string) ([]string, error)
}

type trafficRepository struct {
	db *sqlx.DB
}

func NewTrafficRepository(db *sqlx.DB) *trafficRepository {
	return &trafficRepository{db}
}

func (r *trafficRepository) GetSnmpIndexByODN(ctx context.Context, ip, odn string) ([]string, error) {
	var res []string
	query := `SELECT pd.idx FROM prometheus_devices AS pd
		JOIN fats AS f ON pd.ip = f.ip AND pd.shell = f.shell AND pd.card = f.card AND pd.port = f.port
		WHERE f.ip = $1 AND f.odn = $2;`
	err := r.db.SelectContext(ctx, &res, query, ip, odn)
	return res, err
}

func (r *trafficRepository) GetSnmpIndexByFAT(ctx context.Context, ip, odn, fat string) ([]string, error) {
	var res []string
	query := `SELECT pd.idx FROM prometheus_devices AS pd
		JOIN fats AS f ON pd.ip = f.ip AND pd.shell = f.shell AND pd.card = f.card AND pd.port = f.port
		WHERE f.ip = $1 AND f.odn = $2 AND f.fat = $3;`
	err := r.db.SelectContext(ctx, &res, query, ip, odn, fat)
	return res, err
}
