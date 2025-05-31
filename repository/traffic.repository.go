package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/constants"
)

type TrafficRepository interface {
	TotalTraffic(ctx context.Context, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByCounty(ctx context.Context, state, county string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByMunicipality(ctx context.Context, state, county, municipality string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByPon(ctx context.Context, sysname, ifname string, initDate, endDate time.Time) ([]entity.Traffic, error)
}

type trafficRepository struct {
	db *sqlx.DB
}

func NewTrafficRepository(db *sqlx.DB) *trafficRepository {
	return &trafficRepository{db}
}

func (repo *trafficRepository) TotalTraffic(ctx context.Context, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	err := sqlx.SelectContext(ctx, repo.db, &res, constants.SQL_TOTAL_TRAFFIC, initDate, endDate)
	return res, err
}

func (repo *trafficRepository) TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	err := sqlx.SelectContext(ctx, repo.db, &res, constants.SQL_TRAFFIC_BY_STATE, state, initDate, endDate)
	return res, err
}

func (repo *trafficRepository) TrafficByCounty(ctx context.Context, state, county string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	err := sqlx.SelectContext(ctx, repo.db, &res, constants.SQL_TRAFFIC_BY_COUNTY, state, county, initDate, endDate)
	return res, err
}

func (repo *trafficRepository) TrafficByMunicipality(ctx context.Context, state, county, municipality string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	err := sqlx.SelectContext(ctx, repo.db, &res, constants.SQL_TRAFFIC_BY_MUNICIPALITY, state, county, municipality, initDate, endDate)
	return res, err
}

func (repo *trafficRepository) TrafficByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	err := sqlx.SelectContext(ctx, repo.db, &res, constants.SQL_TRAFFIC_BY_ODN, state, odn, initDate, endDate)
	return res, err
}

func (repo *trafficRepository) TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	err := sqlx.SelectContext(ctx, repo.db, &res, constants.SQL_TRAFFIC_BY_OLT, sysname, initDate, endDate)
	return res, err
}

func (repo *trafficRepository) TrafficByPon(ctx context.Context, sysname, ifname string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	err := repo.db.SelectContext(ctx, &res, constants.SQL_TRAFFIC_BY_PON, sysname, ifname, initDate, endDate)
	return res, err
}
