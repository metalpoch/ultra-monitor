package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type PonRepository interface {
	PonsByOLT(ctx context.Context, sysname string) ([]entity.Pon, error)
	PonByPort(ctx context.Context, sysname, port string) (entity.Pon, error)

	TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByMunicipality(ctx context.Context, state, municipality string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByCounty(ctx context.Context, state, municipality, county string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByODN(ctx context.Context, state, municipality, county, odn string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.Traffic, error)
	TrafficByPon(ctx context.Context, sysname, ifname string, initDate, endDate time.Time) ([]entity.Traffic, error)

	GetDailyAveragedHourlyMaxTrafficTrends(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficSummary, error)
	UpsertSummaryTraffic(ctx context.Context, traffic []entity.TrafficSummary) error
	GetTrafficSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficTotalSummary, error)
	GetTrafficStatesSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error)
	GetTrafficMunicipalitySummary(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error)
	GetTrafficCountySummary(ctx context.Context, state, municipality string, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error)
	GetTrafficOdnSummary(ctx context.Context, state, municipality, county string, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error)
}

type ponRepository struct {
	db *sqlx.DB
}

func NewPonRepository(db *sqlx.DB) *ponRepository {
	return &ponRepository{db}
}

func (repo *ponRepository) PonsByOLT(ctx context.Context, sysname string) ([]entity.Pon, error) {
	var res []entity.Pon
	query := `SELECT pons.* FROM pons JOIN olts ON olts.ip = pons.olt_ip WHERE olts.sys_name = $1`
	err := repo.db.SelectContext(ctx, &res, query, sysname)
	return res, err
}

func (repo *ponRepository) PonByPort(ctx context.Context, sysname, port string) (entity.Pon, error) {
	var res entity.Pon
	query := `SELECT pons.* FROM pons JOIN olts ON olts.ip = pons.olt_ip WHERE olts.sys_name = $1 AND pons.if_name = $2`
	err := repo.db.GetContext(ctx, &res, query, sysname, port)
	return res, err
}

func (repo *ponRepository) TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
    WITH fats_unique AS (
        SELECT *,
            ROW_NUMBER() OVER (PARTITION BY olt_ip ORDER BY id DESC) AS rn
        FROM fats
    ),
    traffic_hourly AS (
        SELECT
            pon_id,
            DATE_TRUNC('hour', date) AS date,
            SUM(bps_in) AS bps_in,
            SUM(bps_out) AS bps_out,
            SUM(bandwidth_mbps_sec) AS bandwidth_mbps_sec,
            SUM(bytes_in) AS bytes_in,
            SUM(bytes_out) AS bytes_out
        FROM traffic_pons
        WHERE date BETWEEN $2 AND $3
        GROUP BY pon_id, date
    )
    SELECT
        th.date,
        SUM(th.bps_in) / 1000000 AS mbps_in,
        SUM(th.bps_out) / 1000000 AS mbps_out,
        SUM(th.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(th.bytes_in) / 1000000 AS mbytes_in,
        SUM(th.bytes_out) / 1000000 AS mbytes_out
    FROM traffic_hourly th
    JOIN pons ON pons.id = th.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    JOIN fats_unique fu ON fu.olt_ip = olts.ip AND fu.rn = 1
    WHERE fu.state = $1
    GROUP BY th.date, fu.state
    ORDER BY th.date;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, initDate, endDate)
	return res, err
}

func (repo *ponRepository) TrafficByMunicipality(ctx context.Context, state, municipality string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
    WITH fats_unique AS (
        SELECT *,
            ROW_NUMBER() OVER (PARTITION BY olt_ip ORDER BY id DESC) AS rn
        FROM fats
    ),
    traffic_hourly AS (
        SELECT
            pon_id,
            DATE_TRUNC('hour', date) AS date,
            SUM(bps_in) AS bps_in,
            SUM(bps_out) AS bps_out,
            SUM(bandwidth_mbps_sec) AS bandwidth_mbps_sec,
            SUM(bytes_in) AS bytes_in,
            SUM(bytes_out) AS bytes_out
        FROM traffic_pons
        WHERE date BETWEEN $3 AND $4
        GROUP BY pon_id, date
    )
    SELECT
        th.date,
        SUM(th.bps_in) / 1000000 AS mbps_in,
        SUM(th.bps_out) / 1000000 AS mbps_out,
        SUM(th.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(th.bytes_in) / 1000000 AS mbytes_in,
        SUM(th.bytes_out) / 1000000 AS mbytes_out
    FROM traffic_hourly th
    JOIN pons ON pons.id = th.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    JOIN fats_unique fu ON fu.olt_ip = olts.ip AND fu.rn = 1
    WHERE fu.state = $1 AND fu.municipality = $2
    GROUP BY th.date, fu.state
    ORDER BY th.date;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, municipality, initDate, endDate)
	return res, err
}

func (repo *ponRepository) TrafficByCounty(ctx context.Context, state, municipality, county string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
    WITH fats_unique AS (
        SELECT *,
            ROW_NUMBER() OVER (PARTITION BY olt_ip ORDER BY id DESC) AS rn
        FROM fats
    ),
    traffic_hourly AS (
        SELECT
            pon_id,
            DATE_TRUNC('hour', date) AS date,
            SUM(bps_in) AS bps_in,
            SUM(bps_out) AS bps_out,
            SUM(bandwidth_mbps_sec) AS bandwidth_mbps_sec,
            SUM(bytes_in) AS bytes_in,
            SUM(bytes_out) AS bytes_out
        FROM traffic_pons
        WHERE date BETWEEN $4 AND $5
        GROUP BY pon_id, date
    )
    SELECT
        th.date,
        SUM(th.bps_in) / 1000000 AS mbps_in,
        SUM(th.bps_out) / 1000000 AS mbps_out,
        SUM(th.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(th.bytes_in) / 1000000 AS mbytes_in,
        SUM(th.bytes_out) / 1000000 AS mbytes_out
    FROM traffic_hourly th
    JOIN pons ON pons.id = th.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    JOIN fats_unique fu ON fu.olt_ip = olts.ip AND fu.rn = 1
    WHERE fu.state = $1 AND fu.municipality = $2  AND fu.county = $3
    GROUP BY th.date, fu.state
    ORDER BY th.date;`

	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, municipality, county, initDate, endDate)
	return res, err
}

func (repo *ponRepository) TrafficByODN(ctx context.Context, state, municipality, county, odn string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
    WITH fats_unique AS (
        SELECT *,
            ROW_NUMBER() OVER (PARTITION BY olt_ip ORDER BY id DESC) AS rn
        FROM fats
    ),
    traffic_hourly AS (
        SELECT
            pon_id,
            DATE_TRUNC('hour', date) AS date,
            SUM(bps_in) AS bps_in,
            SUM(bps_out) AS bps_out,
            SUM(bandwidth_mbps_sec) AS bandwidth_mbps_sec,
            SUM(bytes_in) AS bytes_in,
            SUM(bytes_out) AS bytes_out
        FROM traffic_pons
        WHERE date BETWEEN $5 AND $6
        GROUP BY pon_id, date
    )
    SELECT
        th.date,
        SUM(th.bps_in) / 1000000 AS mbps_in,
        SUM(th.bps_out) / 1000000 AS mbps_out,
        SUM(th.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(th.bytes_in) / 1000000 AS mbytes_in,
        SUM(th.bytes_out) / 1000000 AS mbytes_out
    FROM traffic_hourly th
    JOIN pons ON pons.id = th.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    JOIN fats_unique fu ON fu.olt_ip = olts.ip
    WHERE fu.state = $1 AND fu.municipality = $2  AND fu.county = $3 AND fu.odn = $4
    GROUP BY th.date, fu.state
    ORDER BY th.date;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, municipality, county, odn, initDate, endDate)
	return res, err
}

func (repo *ponRepository) TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
    SELECT
        DATE_TRUNC('hour', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in) / 1000000 AS mbytes_in,
        SUM(traffic_pons.bytes_out) / 1000000 AS mbytes_out
    FROM traffic_pons 
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    WHERE olts.sys_name = $1 AND traffic_pons.date BETWEEN $2 AND $3
    GROUP BY DATE_TRUNC('hour', date)
    ORDER BY date;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, sysname, initDate, endDate)
	return res, err
}

func (repo *ponRepository) TrafficByPon(ctx context.Context, sysname, ifname string, initDate, endDate time.Time) ([]entity.Traffic, error) {
	var res []entity.Traffic
	query := `
    SELECT
        DATE_TRUNC('hour', traffic_pons.date) AS date,
        SUM(traffic_pons.bps_in) / 1000000 AS mbps_in,
        SUM(traffic_pons.bps_out) / 1000000 AS mbps_out,
        SUM(traffic_pons.bandwidth_mbps_sec) / 1000000 AS bandwidth_mbps_sec,
        SUM(traffic_pons.bytes_in) / 1000000 AS mbytes_in,
        SUM(traffic_pons.bytes_out) / 1000000 AS mbytes_out
    FROM traffic_pons 
    JOIN pons ON pons.id = traffic_pons.pon_id
    JOIN olts ON olts.ip = pons.olt_ip
    WHERE olts.sys_name = $1 AND pons.if_name = $2 AND traffic_pons.date BETWEEN $3 AND $4
    GROUP BY DATE_TRUNC('hour', date)
    ORDER BY date;`
	err := repo.db.SelectContext(ctx, &res, query, sysname, ifname, initDate, endDate)
	return res, err
}

func (repo *ponRepository) GetDailyAveragedHourlyMaxTrafficTrends(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficSummary, error) {
	var res []entity.TrafficSummary
	query := `
    WITH hourly_max AS (
        SELECT
            DATE(date) AS day,
            pon_id,
            EXTRACT(HOUR FROM date) AS hour,
            MAX(bps_in) AS max_bps_in,
            MAX(bps_out) AS max_bps_out,
            MAX(bytes_in) AS max_bytes_in,
            MAX(bytes_out) AS max_bytes_out
        FROM traffic_pons
        WHERE date BETWEEN $1 AND $2
        GROUP BY day, pon_id, hour
    ),
    hourly_avg AS (
        SELECT
            day,
            pon_id,
            AVG(max_bps_in) / 1e6 AS mbps_in,
            AVG(max_bps_out) / 1e6 AS mbps_out,
            AVG(max_bytes_in) / 1e6 AS mbytes_in,
            AVG(max_bytes_out) / 1e6 AS mbytes_out
        FROM hourly_max
        GROUP BY day, pon_id
    ),
    joined_data AS (
        SELECT
            hm.day,
            hm.pon_id,
            hm.mbps_in,
            hm.mbps_out,
            hm.mbytes_in,
            hm.mbytes_out,
            p.olt_ip
        FROM hourly_avg hm
        JOIN pons p ON p.id = hm.pon_id
    )
    SELECT
        day,
        olt_ip,
        SUM(mbps_in) AS mbps_in,
        SUM(mbps_out) AS mbps_out,
        SUM(mbytes_in) AS mbytes_in,
        SUM(mbytes_out) AS mbytes_out
    FROM joined_data
    GROUP BY day, olt_ip
    ORDER BY day, olt_ip;`

	err := repo.db.SelectContext(ctx, &res, query, initDate, endDate)
	return res, err
}

func (repo *ponRepository) UpsertSummaryTraffic(ctx context.Context, counts []entity.TrafficSummary) error {
	const fieldCount = 6
	query := `
        INSERT INTO traffic_pons_summary (
            day, olt_ip, mbps_in, mbps_out, mbytes_in, mbytes_out
        ) VALUES `
	valueStrings := make([]string, 0, len(counts))
	valueArgs := make([]interface{}, 0, len(counts)*fieldCount)

	for i, c := range counts {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d)",
			i*fieldCount+1, i*fieldCount+2, i*fieldCount+3, i*fieldCount+4, i*fieldCount+5, i*fieldCount+6))
		valueArgs = append(valueArgs,
			c.Day, c.OltIP, c.MbpsIn, c.MbpsOut, c.MbytesIn, c.MbytesOut)
	}

	query += strings.Join(valueStrings, ", ")
	query += `
	    ON CONFLICT (day, olt_ip) DO UPDATE SET
    	    mbps_in = EXCLUDED.mbps_in,
	        mbps_out = EXCLUDED.mbps_out,
	        mbytes_in = EXCLUDED.mbytes_in,
	        mbytes_out = EXCLUDED.mbytes_out`

	_, err := repo.db.ExecContext(ctx, query, valueArgs...)
	return err
}

func (repo *ponRepository) GetTrafficSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficTotalSummary, error) {
	var res []entity.TrafficTotalSummary
	query := `
    SELECT
        day,
        SUM(mbps_in) AS mbps_in,
        SUM(mbps_out) AS mbps_out,
        SUM(mbytes_in) AS mbytes_in,
        SUM(mbytes_out) AS mbytes_out
    FROM traffic_pons_summary
    WHERE day BETWEEN $1 AND $2
    GROUP BY day
    ORDER BY day;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ponRepository) GetTrafficStatesSummary(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error) {
	var res []entity.TrafficInfoSummary
	query := `
    WITH unique_traffic AS (
        SELECT
            day,
            olt_ip,
            SUM(mbps_in) AS mbps_in,
            SUM(mbps_out) AS mbps_out,
            SUM(mbytes_in) AS mbytes_in,
            SUM(mbytes_out) AS mbytes_out
        FROM traffic_pons_summary
        WHERE day BETWEEN $1 AND $2
        GROUP BY day, olt_ip
    ),
    unique_fats AS (
        SELECT DISTINCT ON (olt_ip) olt_ip, state FROM fats
    )
    SELECT
        ut.day,
        f.state AS description,
        SUM(ut.mbps_in) AS mbps_in,
        SUM(ut.mbps_out) AS mbps_out,
        SUM(ut.mbytes_in) AS mbytes_in,
        SUM(ut.mbytes_out) AS mbytes_out
    FROM unique_traffic ut
    JOIN unique_fats f ON f.olt_ip = ut.olt_ip
    GROUP BY ut.day, f.state
    ORDER BY ut.day, f.state;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	return res, err
}

func (repo *ponRepository) GetTrafficMunicipalitySummary(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error) {
	var res []entity.TrafficInfoSummary
	query := `
    WITH unique_traffic AS (
        SELECT
            day,
            olt_ip,
            SUM(mbps_in) AS mbps_in,
            SUM(mbps_out) AS mbps_out,
            SUM(mbytes_in) AS mbytes_in,
            SUM(mbytes_out) AS mbytes_out
        FROM traffic_pons_summary
        WHERE day BETWEEN $1 AND $2
        GROUP BY day, olt_ip
    ),
    unique_fats AS (
        SELECT DISTINCT ON (olt_ip) olt_ip, municipality FROM fats WHERE state = $3
    )
    SELECT
        ut.day,
        f.municipality AS description,
        SUM(ut.mbps_in) AS mbps_in,
        SUM(ut.mbps_out) AS mbps_out,
        SUM(ut.mbytes_in) AS mbytes_in,
        SUM(ut.mbytes_out) AS mbytes_out
    FROM unique_traffic ut
    JOIN unique_fats f ON f.olt_ip = ut.olt_ip
    GROUP BY ut.day, f.municipality
    ORDER BY ut.day, f.municipality;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"), state)
	return res, err
}

func (repo *ponRepository) GetTrafficCountySummary(ctx context.Context, state, municipality string, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error) {
	var res []entity.TrafficInfoSummary
	query := `
    WITH unique_traffic AS (
        SELECT
            day,
            olt_ip,
            SUM(mbps_in) AS mbps_in,
            SUM(mbps_out) AS mbps_out,
            SUM(mbytes_in) AS mbytes_in,
            SUM(mbytes_out) AS mbytes_out
        FROM traffic_pons_summary
        WHERE day BETWEEN $1 AND $2
        GROUP BY day, olt_ip
    ),
    unique_fats AS (
        SELECT DISTINCT ON (olt_ip) olt_ip, county FROM fats WHERE state = $3 AND municipality = $4
    )
    SELECT
        ut.day,
        f.county AS description,
        SUM(ut.mbps_in) AS mbps_in,
        SUM(ut.mbps_out) AS mbps_out,
        SUM(ut.mbytes_in) AS mbytes_in,
        SUM(ut.mbytes_out) AS mbytes_out
    FROM unique_traffic ut
    JOIN unique_fats f ON f.olt_ip = ut.olt_ip
    GROUP BY ut.day, f.county
    ORDER BY ut.day, f.county;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"), state, municipality)
	return res, err
}

func (repo *ponRepository) GetTrafficOdnSummary(ctx context.Context, state, municipality, county string, initDate, endDate time.Time) ([]entity.TrafficInfoSummary, error) {
	var res []entity.TrafficInfoSummary
	query := `
    WITH unique_traffic AS (
        SELECT
            day,
            olt_ip,
            SUM(mbps_in) AS mbps_in,
            SUM(mbps_out) AS mbps_out,
            SUM(mbytes_in) AS mbytes_in,
            SUM(mbytes_out) AS mbytes_out
        FROM traffic_pons_summary
        WHERE day BETWEEN $1 AND $2
        GROUP BY day, olt_ip
    ),
    unique_fats AS (
        SELECT DISTINCT ON (olt_ip) olt_ip, odn FROM fats WHERE state = $3 AND municipality = $4 AND county = $5
    )
    SELECT
        ut.day,
        f.odn AS description,
        SUM(ut.mbps_in) AS mbps_in,
        SUM(ut.mbps_out) AS mbps_out,
        SUM(ut.mbytes_in) AS mbytes_in,
        SUM(ut.mbytes_out) AS mbytes_out
    FROM unique_traffic ut
    JOIN unique_fats f ON f.olt_ip = ut.olt_ip
    GROUP BY ut.day, f.odn
    ORDER BY ut.day, f.odn;`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, initDate.Format("2006-01-02"), endDate.Format("2006-01-02"), state, municipality, county)
	return res, err
}
