package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type OntRepository interface {
	GetOntStatusByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error)
	GetOntStatusByODN(ctx context.Context, state, municipality, county, odn string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error)
	GetOntStatusByOltIP(ctx context.Context, ip string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error)
	GetOntStatusBySysname(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error)
	TrafficOnt(ctx context.Context, ponID int, idx int64, initDate, endDate time.Time) ([]entity.TrafficOnt, error)
	TrafficOntByDespt(ctx context.Context, despt string, initDate, endDate time.Time) ([]entity.TrafficOnt, error)
}

type ontRepository struct {
	db *sqlx.DB
}

func NewOntRepository(db *sqlx.DB) *ontRepository {
	return &ontRepository{db}
}

func (repo *ontRepository) GetOntStatusByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error) {
	var res []entity.OntStatusCountsByState
	query := `
	WITH ranked_status AS (
        SELECT
            olts.sys_name AS sysname,
            DATE_TRUNC('day', measurement_onts.date) AS date,
            measurement_onts.pon_id,
            measurement_onts.idx,
            MIN(
                CASE
                    WHEN control_run_status = 1 THEN 1
                    WHEN control_run_status = 2 THEN 2
                    ELSE 3
                END
            ) AS status_priority
        FROM measurement_onts
        JOIN pons ON measurement_onts.pon_id = pons.id
        JOIN olts ON pons.olt_ip = olts.ip
        JOIN fats ON fats.olt_ip = olts.ip
        WHERE fats.state = $1 AND measurement_onts.date BETWEEN $2 AND $3
        GROUP BY sysname, DATE_TRUNC('day', measurement_onts.date), measurement_onts.pon_id, measurement_onts.idx
    )
    SELECT
        sysname,
        date,
        COUNT(DISTINCT pon_id) AS ports_pon,
        SUM(CASE WHEN status_priority = 1 THEN 1 ELSE 0 END) AS actives,
        SUM(CASE WHEN status_priority = 2 THEN 1 ELSE 0 END) AS inactives,
        SUM(CASE WHEN status_priority = 3 THEN 1 ELSE 0 END) AS unknowns,
        COUNT(*) AS total
    FROM ranked_status
    GROUP BY sysname, date
    ORDER BY sysname, date;`
	err := repo.db.SelectContext(ctx, &res, query, state, initDate, endDate)
	return res, err
}

func (repo *ontRepository) GetOntStatusByODN(ctx context.Context, state, municipality, county, odn string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error) {
	var res []entity.OntStatusCountsByState
	query := `
	WITH ranked_status AS (
        SELECT
            olts.sys_name AS sysname,
            DATE_TRUNC('day', measurement_onts.date) AS date,
            measurement_onts.pon_id,
            measurement_onts.idx,
            MIN(
                CASE
                    WHEN control_run_status = 1 THEN 1
                    WHEN control_run_status = 2 THEN 2
                    ELSE 3
                END
            ) AS status_priority
        FROM measurement_onts
        JOIN pons ON measurement_onts.pon_id = pons.id
        JOIN olts ON pons.olt_ip = olts.ip
        JOIN fats ON fats.olt_ip = olts.ip
		WHERE 
			fats.state = $1
			AND fats.municipality = $2
			AND fats.county = $3
			AND fats.odn = $4
			AND measurement_onts.date BETWEEN $5 AND $6
        GROUP BY sysname, DATE_TRUNC('day', measurement_onts.date), measurement_onts.pon_id, measurement_onts.idx
    )
    SELECT
        sysname,
        date,
        COUNT(DISTINCT pon_id) AS ports_pon,
        SUM(CASE WHEN status_priority = 1 THEN 1 ELSE 0 END) AS actives,
        SUM(CASE WHEN status_priority = 2 THEN 1 ELSE 0 END) AS inactives,
        SUM(CASE WHEN status_priority = 3 THEN 1 ELSE 0 END) AS unknowns,
        COUNT(*) AS total
    FROM ranked_status
    GROUP BY sysname, date
    ORDER BY sysname, date;`

	err := repo.db.SelectContext(ctx, &res, query, state, municipality, county, odn, initDate, endDate)
	return res, err
}

func (repo *ontRepository) GetOntStatusByOltIP(ctx context.Context, ip string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error) {
	var res []entity.OntStatusCountsByState
	query := `
	WITH ranked_status AS (
        SELECT
            olts.sys_name AS sysname,
            DATE_TRUNC('day', measurement_onts.date) AS date,
            measurement_onts.pon_id,
            measurement_onts.idx,
            MIN(
                CASE
                    WHEN control_run_status = 1 THEN 1
                    WHEN control_run_status = 2 THEN 2
                    ELSE 3
                END
            ) AS status_priority
        FROM measurement_onts
        JOIN pons ON measurement_onts.pon_id = pons.id
        JOIN olts ON pons.olt_ip = olts.ip
        WHERE olts.ip = $1 AND measurement_onts.date BETWEEN $2 AND $3
        GROUP BY sysname, DATE_TRUNC('day', measurement_onts.date), measurement_onts.pon_id, measurement_onts.idx
    )
    SELECT
        sysname,
        date,
        COUNT(DISTINCT pon_id) AS ports_pon,
        SUM(CASE WHEN status_priority = 1 THEN 1 ELSE 0 END) AS actives,
        SUM(CASE WHEN status_priority = 2 THEN 1 ELSE 0 END) AS inactives,
        SUM(CASE WHEN status_priority = 3 THEN 1 ELSE 0 END) AS unknowns,
        COUNT(*) AS total
    FROM ranked_status
    GROUP BY sysname, date
    ORDER BY sysname, date;`
	err := repo.db.SelectContext(ctx, &res, query, ip, initDate, endDate)
	return res, err
}

func (repo *ontRepository) GetOntStatusBySysname(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error) {
	var res []entity.OntStatusCountsByState
	query := `
	WITH ranked_status AS (
        SELECT
            olts.sys_name AS sysname,
            DATE_TRUNC('day', measurement_onts.date) AS date,
            measurement_onts.pon_id,
            measurement_onts.idx,
            MIN(
                CASE
                    WHEN control_run_status = 1 THEN 1
                    WHEN control_run_status = 2 THEN 2
                    ELSE 3
                END
            ) AS status_priority
        FROM measurement_onts
        JOIN pons ON measurement_onts.pon_id = pons.id
        JOIN olts ON pons.olt_ip = olts.ip
        WHERE olts.sys_name = $1 AND measurement_onts.date BETWEEN $2 AND $3
        GROUP BY sysname, DATE_TRUNC('day', measurement_onts.date), measurement_onts.pon_id, measurement_onts.idx
    )
    SELECT
        sysname,
        date,
        COUNT(DISTINCT pon_id) AS ports_pon,
        SUM(CASE WHEN status_priority = 1 THEN 1 ELSE 0 END) AS actives,
        SUM(CASE WHEN status_priority = 2 THEN 1 ELSE 0 END) AS inactives,
        SUM(CASE WHEN status_priority = 3 THEN 1 ELSE 0 END) AS unknowns,
        COUNT(*) AS total
    FROM ranked_status
    GROUP BY sysname, date
    ORDER BY sysname, date;`

	err := repo.db.SelectContext(ctx, &res, query, sysname, initDate, endDate)
	return res, err
}

func (repo *ontRepository) TrafficOnt(ctx context.Context, ponID int, idx int64, initDate, endDate time.Time) ([]entity.TrafficOnt, error) {
	var res []entity.TrafficOnt
	query := `
	SELECT
        date,
        despt,
        serial_number,
        line_prof_name,
        olt_distance,
        control_mac_count,
        control_run_status,
        CASE
            WHEN curr_bytes_in < prev_bytes_in THEN ((18446744073709551615 - prev_bytes_in) + curr_bytes_in) * 8 / (time_diff * 1000000)
            ELSE ((curr_bytes_in - prev_bytes_in) * 8) / (time_diff * 1000000)
        END AS mbps_in,
        CASE
            WHEN curr_bytes_out < prev_bytes_out THEN ((18446744073709551615 - prev_bytes_out) + curr_bytes_out) * 8 / (time_diff * 1000000)
            ELSE ((curr_bytes_out - prev_bytes_out) * 8) / (time_diff * 1000000)
        END AS mbps_out,
        CASE
            WHEN curr_bytes_in < prev_bytes_in THEN ((18446744073709551615 - prev_bytes_in) + curr_bytes_in) / (time_diff * 1000000)
            ELSE (curr_bytes_in - prev_bytes_in) / (time_diff * 1000000)
        END AS mbytes_in_sec,
        CASE
            WHEN curr_bytes_out < prev_bytes_out THEN ((18446744073709551615 - prev_bytes_out) + curr_bytes_out) / (time_diff * 1000000)
            ELSE (curr_bytes_out - prev_bytes_out) / (time_diff * 1000000)
        END AS mbytes_out_sec
    FROM (
        SELECT
            date,
            despt,
            serial_number,
            line_prof_name,
            olt_distance,
            control_mac_count,
            control_run_status,
            bytes_in_count AS prev_bytes_in,
            bytes_out_count AS prev_bytes_out,
            LEAD(bytes_in_count) OVER (PARTITION BY pon_id ORDER BY date) AS curr_bytes_in,
            LEAD(bytes_out_count) OVER (PARTITION BY pon_id ORDER BY date) AS curr_bytes_out,
            EXTRACT(EPOCH FROM (LEAD(date) OVER (PARTITION BY pon_id ORDER BY date) - date)) AS time_diff
        FROM measurement_onts
        WHERE pon_id = $1 AND idx = $2 AND bytes_in_count > 0 AND bytes_out_count > 0 AND date BETWEEN $3 AND $4
        ORDER BY date
    ) AS subquery
    WHERE curr_bytes_in IS NOT NULL
      AND curr_bytes_out IS NOT NULL
      AND time_diff IS NOT NULL;`

	err := repo.db.SelectContext(ctx, &res, query, ponID, idx, initDate, endDate)
	return res, err
}

func (repo *ontRepository) TrafficOntByDespt(ctx context.Context, despt string, initDate, endDate time.Time) ([]entity.TrafficOnt, error) {
	var res []entity.TrafficOnt
	query := `SELECT
        date,
        despt,
        serial_number,
        line_prof_name,
        olt_distance,
        control_mac_count,
        control_run_status,
        CASE
            WHEN curr_bytes_in < prev_bytes_in THEN ((18446744073709551615 - prev_bytes_in) + curr_bytes_in) * 8 / (time_diff * 1000000)
            ELSE ((curr_bytes_in - prev_bytes_in) * 8) / (time_diff * 1000000)
        END AS mbps_in,
        CASE
            WHEN curr_bytes_out < prev_bytes_out THEN ((18446744073709551615 - prev_bytes_out) + curr_bytes_out) * 8 / (time_diff * 1000000)
            ELSE ((curr_bytes_out - prev_bytes_out) * 8) / (time_diff * 1000000)
        END AS mbps_out,
        CASE
            WHEN curr_bytes_in < prev_bytes_in THEN ((18446744073709551615 - prev_bytes_in) + curr_bytes_in) / (time_diff * 1000000)
            ELSE (curr_bytes_in - prev_bytes_in) / (time_diff * 1000000)
        END AS mbytes_in_sec,
        CASE
            WHEN curr_bytes_out < prev_bytes_out THEN ((18446744073709551615 - prev_bytes_out) + curr_bytes_out) / (time_diff * 1000000)
            ELSE (curr_bytes_out - prev_bytes_out) / (time_diff * 1000000)
        END AS mbytes_out_sec
    FROM (
        SELECT
            date,
            despt,
            serial_number,
            line_prof_name,
            olt_distance,
            control_mac_count,
            control_run_status,
            bytes_in_count AS prev_bytes_in,
            bytes_out_count AS prev_bytes_out,
            LEAD(bytes_in_count) OVER (PARTITION BY pon_id ORDER BY date) AS curr_bytes_in,
            LEAD(bytes_out_count) OVER (PARTITION BY pon_id ORDER BY date) AS curr_bytes_out,
            EXTRACT(EPOCH FROM (LEAD(date) OVER (PARTITION BY pon_id ORDER BY date) - date)) AS time_diff
        FROM measurement_onts
        WHERE despt = $1 AND bytes_in_count > 0 AND bytes_out_count > 0 AND date BETWEEN $2 AND $3
        ORDER BY date
    ) AS subquery
    WHERE curr_bytes_in IS NOT NULL
      AND curr_bytes_out IS NOT NULL
      AND time_diff IS NOT NULL;`

	err := repo.db.SelectContext(ctx, &res, query, despt, initDate, endDate)
	return res, err
}
