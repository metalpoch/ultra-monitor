package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/constants"
)

type OntRepository interface {
	AllOntStatus(ctx context.Context, initDate, endDate time.Time) ([]entity.OntStatusCounts, error)
	GetOntStatusByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error)
	GetOntStatusByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error)
	TrafficOnt(ctx context.Context, ponID uint64, idx string, initDate, endDate time.Time) ([]entity.TrafficOnt, error)
}

type ontRepository struct {
	db *sqlx.DB
}

func NewOntRepository(db *sqlx.DB) *ontRepository {
	return &ontRepository{db}
}

func (repo *ontRepository) AllOntStatus(ctx context.Context, initDate, endDate time.Time) ([]entity.OntStatusCounts, error) {
	var res []entity.OntStatusCounts
	err := repo.db.SelectContext(ctx, &res, constants.SQL_ALL_ONT_STATUS, initDate, endDate)
	return res, err
}

func (repo *ontRepository) GetOntStatusByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error) {
	var res []entity.OntStatusCountsByState
	err := repo.db.SelectContext(ctx, &res, constants.SQL_ONT_STATUS_BY_STATE, state, initDate, endDate)
	return res, err
}

func (repo *ontRepository) GetOntStatusByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error) {
	var res []entity.OntStatusCountsByState
	err := repo.db.SelectContext(ctx, &res, constants.SQL_ONT_STATUS_BY_ODN, state, odn, initDate, endDate)
	return res, err
}

func (repo *ontRepository) TrafficOnt(ctx context.Context, PonID uint64, idx string, initDate, endDate time.Time) ([]entity.TrafficOnt, error) {
	var res []entity.TrafficOnt
	err := repo.db.SelectContext(ctx, &res, constants.SQL_TRAFFIC_ONT, PonID, idx, initDate, endDate)
	return res, err
}
