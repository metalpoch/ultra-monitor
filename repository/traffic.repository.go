package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/entity"
)

type TrafficRepository interface {
	TotalTraffic(ctx context.Context, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByCounty(ctx context.Context, state, county string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByMunicipality(ctx context.Context, state, county, municipality string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByPON(ctx context.Context, sysname, pon string, initDate, endDate time.Time) ([]entity.Traffic, error)
}

type trafficRepository struct {
	db *sqlx.DB
}

func NewTrafficRepository(db *sqlx.DB) *trafficRepository {
	return &trafficRepository{db}
}

func (repo trafficRepository) Traffic(ctx context.Context, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth_mbps_sec,
			SUM(bytes_in) / 1000000 AS mbytes_in_sec,
			SUM(bytes_out) / 1000000 AS mbytes_out_sec
		FROM traffic_pon
		WHERE date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, initDate, endDate)
	return res, err
}

func (repo trafficRepository) TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffic_pon
		JOIN fats_pon ON fats_pon.interface_id = traffic_pon.interface_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON locations.id = fats.location_id
		WHERE locations.state = ? AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, initDate, endDate)
	return res, err
}

func (repo trafficRepository) TrafficByCounty(ctx context.Context, state, county string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffic_pon
		JOIN fats_pon ON fats_pon.interface_id = traffic_pon.interface_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON locations.id = fats.location_id
		WHERE locations.state = ? AND locations.county AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, county, initDate, endDate)
	return res, err
}

func (repo trafficRepository) TrafficByMunicipality(ctx context.Context, state, county, municipality string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffic_pon
		JOIN fats_pon ON fats_pon.interface_id = traffic_pon.interface_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON locations.id = fats.location_id
		WHERE locations.state = ? AND locations.county AND locations.municipality AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, county, municipality, initDate, endDate)
	return res, err
}

func (repo trafficRepository) TrafficByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffic_pon
		JOIN fats_pon ON fats_pon.interface_id = traffic_pon.interface_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON locations.id = fats.location_id
		WHERE locations.state = ? AND fats.odn = ? AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, odn, initDate, endDate)
	return res, err
}

func (repo trafficRepository) TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffic_pon
		JOIN pons ON pons.id = traffic_pon.interface_id"
		JOIN devices ON devices.id = pons.device_id
		WHERE devices.sys_name = ? AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, sysname, initDate, endDate)
	return res, err
}

func (repo trafficRepository) TrafficByPON(ctx context.Context, sysname, pon string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffic_pon
		JOIN pons ON pons.id = traffic_pon.interface_id"
		JOIN devices ON devices.id = pons.device_id
		WHERE devices.sys_name = ? AND pons.if_name = ? AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := repo.db.SelectContext(ctx, &res, query, sysname, pon, initDate, endDate)
	return res, err
}
