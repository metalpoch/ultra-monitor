package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type TrafficRepository interface {
	GetSnmpIndexByMunicipality(ctx context.Context, state, municipality string) ([]entity.OltIndex, error)
	GetSnmpIndexByCounty(ctx context.Context, state, municipality, county string) ([]entity.OltIndex, error)
	GetSnmpIndexByODN(ctx context.Context, state, municipality, odn string) ([]entity.OltIndex, error)
}

type trafficRepository struct {
	db *sqlx.DB
}

func NewTrafficRepository(db *sqlx.DB) *trafficRepository {
	return &trafficRepository{db}
}

func (r *trafficRepository) GetSnmpIndexByMunicipality(ctx context.Context, state, municipality string) ([]entity.OltIndex, error) {
	var res []entity.OltIndex
	query := `SELECT pd.ip, pd.idx FROM prometheus_devices AS pd
		JOIN fats AS f ON pd.ip = f.ip AND pd.shell = f.shell AND pd.card = f.card AND pd.port = f.port
		WHERE f.state = $1 AND f.municipality = $2;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality)
	return res, err
}

func (r *trafficRepository) GetSnmpIndexByCounty(ctx context.Context, state, municipality, county string) ([]entity.OltIndex, error) {
	var res []entity.OltIndex
	query := `SELECT pd.ip, pd.idx FROM prometheus_devices AS pd
		JOIN fats AS f ON pd.ip = f.ip AND pd.shell = f.shell AND pd.card = f.card AND pd.port = f.port
		WHERE f.state = $1 AND f.municipality = $2 AND f.county = $3;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county)
	return res, err
}

func (r *trafficRepository) GetSnmpIndexByODN(ctx context.Context, state, municipality, odn string) ([]entity.OltIndex, error) {
	var res []entity.OltIndex
	query := `SELECT pd.ip, pd.idx FROM prometheus_devices AS pd
		JOIN fats AS f ON pd.ip = f.ip AND pd.shell = f.shell AND pd.card = f.card AND pd.port = f.port
		WHERE f.state = $1 AND f.municipality = $2 AND f.odn = $3;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, odn)
	return res, err
}
