package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
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
	query := `
		SELECT
			fats.state AS state,
			DATE_TRUNC('date', measurement_onts.date) AS date,
    		COUNT(DISTINCT pons.id) AS pons_count,
    		COUNT(CASE WHEN measurement_onts.control_run_status = 1 THEN 1 END) AS active_count,
    		COUNT(CASE WHEN measurement_onts.control_run_status = 2 THEN 1 END) AS inactive_count,
    		COUNT(CASE WHEN measurement_onts.control_run_status NOT IN (1,2) THEN 1 END) AS unknown_count,
    		COUNT(*) AS total_count
		FROM measurement_onts
		JOIN pons ON measurement_onts.pon_id = pons.id
		JOIN olts ON pons.olt_id = olts.id
		JOIN fats ON fats.olt_ip = olts.ip
		WHERE measurement_onts.date BETWEEN $2 AND $3
		GROUP BY state, date
		ORDER BY state, date`
	err := repo.db.SelectContext(ctx, &res, query, initDate, endDate)
	return res, err
}

func (repo *ontRepository) GetOntStatusByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error) {
	var res []entity.OntStatusCountsByState
	query := `
		SELECT
			olts.sys_name AS sysname,
			DATE_TRUNC('date', measurement_onts.date) AS date,
    		COUNT(DISTINCT pons.id) AS pons_count,
    		COUNT(CASE WHEN measurement_onts.control_run_status = 1 THEN 1 END) AS active_count,
    		COUNT(CASE WHEN measurement_onts.control_run_status = 2 THEN 1 END) AS inactive_count,
    		COUNT(CASE WHEN measurement_onts.control_run_status NOT IN (1,2) THEN 1 END) AS unknown_count,
    		COUNT(*) AS total_count
		FROM measurement_onts
		JOIN pons ON measurement_onts.pon_id = pons.id
		JOIN olts ON pons.olt_id = olts.id
		JOIN fats ON fats.olt_ip = olts.ip
		WHERE fats.state = $1 AND measurement_onts.date BETWEEN $2 AND $3
		GROUP BY sysname, date
		ORDER BY sysname, date`
	err := repo.db.SelectContext(ctx, &res, query, state, initDate, endDate)
	return res, err
}

func (repo *ontRepository) GetOntStatusByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error) {
	var res []entity.OntStatusCountsByState
	query := `
		SELECT
			olts.sys_name AS sysname,
			DATE_TRUNC('date', measurement_onts.date) AS date,
    		COUNT(DISTINCT pons.id) AS pons_count,
    		COUNT(CASE WHEN measurement_onts.control_run_status = 1 THEN 1 END) AS active_count,
    		COUNT(CASE WHEN measurement_onts.control_run_status = 2 THEN 1 END) AS inactive_count,
    		COUNT(CASE WHEN measurement_onts.control_run_status NOT IN (1,2) THEN 1 END) AS unknown_count,
    		COUNT(*) AS total_count
		FROM measurement_onts
		JOIN pons ON measurement_onts.pon_id = pons.id
		JOIN olts ON pons.olt_id = olts.id
		JOIN fats ON fats.olt_ip = olts.ip
		WHERE fats.state = $1 AND fats.odn = $2 AND measurement_onts.date BETWEEN $3 AND $4
		GROUP BY sysname, date
		ORDER BY sysname, date`
	err := repo.db.SelectContext(ctx, &res, query, state, odn, initDate, endDate)
	return res, err
}

func (repo *ontRepository) TrafficOnt(ctx context.Context, PonID uint64, idx string, initDate, endDate time.Time) ([]entity.TrafficOnt, error) {
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
			END AS Mbps_in,
			CASE
			WHEN curr_bytes_out < prev_bytes_out THEN ((18446744073709551615 - prev_bytes_out) + curr_bytes_out) * 8 / (time_diff * 1000000)
			ELSE ((curr_bytes_out - prev_bytes_out) * 8) / (time_diff * 1000000)
			END AS Mbps_out,
			CASE
			WHEN curr_bytes_in < prev_bytes_in THEN ((18446744073709551615 - prev_bytes_in) + curr_bytes_in) / (time_diff * 1000000)
			ELSE (curr_bytes_in - prev_bytes_in) / (time_diff * 1000000)
			END AS Mbytes_in_sec,
			CASE
			WHEN curr_bytes_out < prev_bytes_out THEN ((18446744073709551615 - prev_bytes_out) + curr_bytes_out) / (time_diff * 1000000)
			ELSE (curr_bytes_out - prev_bytes_out) / (time_diff * 1000000)
			END AS Mbytes_out_sec
		FROM (
			SELECT
				date,
				despt,
				serial_number,
        			line_prof_name,
				olt_distance,
        			control_mac_count, 
				control_run_status,
				bytes_in AS prev_bytes_in,
				bytes_out AS prev_bytes_out,
				LEAD(bytes_in) OVER (PARTITION BY pon_id ORDER BY date) AS curr_bytes_in,
				LEAD(bytes_out) OVER (PARTITION BY pon_id ORDER BY date) AS curr_bytes_out,
				EXTRACT(EPOCH FROM (LEAD(date) OVER (PARTITION BY pon_id ORDER BY date) - date)) AS time_diff
			FROM measurement_ont
			WHERE pon_id = $1 AND idx = $2 AND bytes_in > 0 AND bytes_out > 0 AND date BETWEEN $3 AND $4
			ORDER BY date
		) AS subquery`

	err := repo.db.SelectContext(ctx, &res, query, PonID, idx, initDate, endDate)
	return res, err
}
