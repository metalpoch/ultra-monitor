package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
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
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth_mbps_sec,
			SUM(bytes_in) / 1000000 AS mbytes_in_sec,
			SUM(bytes_out) / 1000000 AS mbytes_out_sec
		FROM traffic_pons
		WHERE date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, initDate, endDate)
	return res, err
}

func (repo *trafficRepository) TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
        SELECT
            DATE_TRUNC('minute', traffic_pons.date) AS date,
            SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
            SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
            SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
            SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
            SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
        FROM traffic_pons
        JOIN pons ON pons.id = traffic_pons.pon_id
        JOIN olts ON olts.id = pons.olt_id
        JOIN fats ON fats.olt_ip = olts.ip
        WHERE fats.state = $1 AND traffic_pons.date BETWEEN $2 AND $3
        GROUP BY DATE_TRUNC('minute', date)
        ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, initDate, endDate)
	return res, err
}

func (repo *trafficRepository) TrafficByCounty(ctx context.Context, state, county string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
        SELECT
            DATE_TRUNC('minute', traffic_pons.date) AS date,
            SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
            SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
            SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
            SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
            SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
        FROM traffic_pons
        JOIN pons ON pons.id = traffic_pons.pon_id
        JOIN olts ON olts.id = pons.olt_id
        JOIN fats ON fats.olt_ip = olts.ip
        WHERE fats.state = $1 AND fats.county = $2 AND traffic_pons.date BETWEEN $3 AND $4
        GROUP BY DATE_TRUNC('minute', date)
        ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, county, initDate, endDate)
	return res, err
}

func (repo *trafficRepository) TrafficByMunicipality(ctx context.Context, state, county, municipality string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
         SELECT
            DATE_TRUNC('minute', traffic_pons.date) AS date,
            SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
            SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
            SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
            SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
            SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
        FROM traffic_pons
        JOIN pons ON pons.id = traffic_pons.pon_id
        JOIN olts ON olts.id = pons.olt_id
        JOIN fats ON fats.olt_ip = olts.ip
        WHERE fats.state = $1 AND fats.county = $2 AND fats.municipality = $3 AND traffic_pons.date BETWEEN $4 AND $5
        GROUP BY DATE_TRUNC('minute', date)
        ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, county, municipality, initDate, endDate)
	return res, err
}

func (repo *trafficRepository) TrafficByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
         SELECT
            DATE_TRUNC('minute', traffic_pons.date) AS date,
            SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
            SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
            SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
            SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
            SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
        FROM traffic_pons
        JOIN pons ON pons.id = traffic_pons.pon_id
        JOIN olts ON olts.id = pons.olt_id
        JOIN fats ON fats.olt_ip = olts.ip
        WHERE fats.state = $1 AND fats.odn = $2 AND traffic_pons.date BETWEEN $3 AND $4
        GROUP BY DATE_TRUNC('minute', date)
        ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, odn, initDate, endDate)
	return res, err
}

func (repo *trafficRepository) TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
        SELECT
            DATE_TRUNC('minute', traffic_pons.date) AS date,
            SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
            SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
            SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
            SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
            SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
        FROM traffic_pons 
        JOIN pons ON pons.id = traffic_pons.pon_id
        JOIN olts ON olts.id = pons.olt_id
        WHERE olts.sys_name = $1 AND traffic_pons.date BETWEEN $2 AND $3
        GROUP BY DATE_TRUNC('minute', date)
        ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, sysname, initDate, endDate)
	return res, err
}

func (repo *trafficRepository) TrafficByPon(ctx context.Context, sysname, ifname string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
        SELECT
            DATE_TRUNC('minute', traffic_pons.date) AS date,
            SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
            SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
            SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
            SUM(traffic_pons.bytes_in_sec) / 1000000 AS mbytes_in,
            SUM(traffic_pons.bytes_out_sec) / 1000000 AS mbytes_out
        FROM traffic_pons 
        JOIN pons ON pons.id = traffic_pons.pon_id
        JOIN olts ON olts.id = pons.olt_id
        WHERE olts.sys_name = $1 AND pons.if_name = $2 AND traffic_pons.date BETWEEN $3 AND $4
        GROUP BY DATE_TRUNC('minute', date)
        ORDER BY date`
	err := repo.db.SelectContext(ctx, &res, query, sysname, ifname, initDate, endDate)
	return res, err
}
