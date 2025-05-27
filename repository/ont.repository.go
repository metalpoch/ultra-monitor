package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/entity"
)

type OntRepository interface {
	AddTrafficOnt(ctx context.Context, traffic *entity.TrafficOnt) error
	GetOntStatus(ctx context.Context, initDate, endDate time.Time) ([]entity.OntStatusCounts, error)
	GetOntStatusByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error)
	GetOntStatusByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error)
	GetTrafficOnt(ctx context.Context, interfaceID, idx string, initDate, endDate time.Time) ([]entity.TrafficOnt, error)
}

type ontRepository struct {
	db *sqlx.DB
}

func NewOntRepository(db *sqlx.DB) *ontRepository {
	return &ontRepository{db}
}

func (repo oltRepository) AddTrafficOnt(ctx context.Context, traffic *entity.TrafficOnt) error {
	query := `
	INSERT INTO traffic_olt (
		date,
		despt,
		serial_number,
		olt_distance,
		control_mac_count,
		control_run_status,
		line_prof_name,
		mbps_in,
		mbps_out,
		mbytes_in_sec,
		mbytes_out_sec
	)
	VALUES (
		:date,
		:despt,
		:serial_number,
		:olt_distance,
		:control_mac_count,
		:control_run_status,
		:line_prof_name,
		:mbps_in,
		:mbps_out,
		:mbytes_in_sec,
		:mbytes_out_sec
	)`
	_, err := repo.db.NamedExecContext(ctx, query, traffic)
	return err
}

func (repo ontRepository) GetOntStatus(ctx context.Context, initDate, endDate time.Time) ([]entity.OntStatusCounts, error) {
	var res []entity.OntStatusCounts
	query := `
		SELECT
        	l.state,
        	DATE_TRUNC('hour', m.date) AS hour,
        	COUNT(DISTINCT m.interface_id) AS pons_count,
        	COUNT(CASE WHEN control_run_status = 1 THEN 1 END) AS active_count,
        	COUNT(CASE WHEN control_run_status = 2 THEN 1 END) AS inactive_count,
        	COUNT(CASE WHEN control_run_status NOT IN (1,2) THEN 1 END) AS unknown_count,
        	COUNT(*) AS total_count
    	FROM measurement_ont AS m
    	JOIN fats_pon ON fats_pon.interface_id = m.interface_id
    	JOIN fats ON fats.id = fats_pon.fat_id
    	JOIN locations AS l ON l.id = fats.location_id
    	WHERE m.date >= $1 AND m.date < $2
    	GROUP BY l.state, hour
    	ORDER BY l.state, hour`
	err := repo.db.SelectContext(ctx, &res, query, initDate, endDate)
	return res, err
}

func (repo ontRepository) GetOntStatusByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error) {
	var res []entity.OntStatusCountsByState
	query := `
		SELECT
			devices.sys_name AS sysname,
			DATE_TRUNC('hour', measurement_ont.date) AS hour,
			COUNT(DISTINCT measurement_ont.interface_id) AS pons_count,
			COUNT(CASE WHEN control_run_status = 1 THEN 1 END) AS active_count,
			COUNT(CASE WHEN control_run_status = 2 THEN 1 END) AS inactive_count,
			COUNT(CASE WHEN control_run_status NOT IN (1,2) THEN 1 END) AS unknown_count,
			COUNT(*) AS total_count
		FROM measurement_ont
		INNER JOIN fats_pon ON fats_pon.interface_id = measurement_ont.interface_id
		INNER JOIN fats ON fats.id = fats_pon.fat_id
		INNER JOIN locations ON locations.id = fats.location_id
		INNER JOIN pons ON pons.id = measurement_ont.interface_id
		INNER JOIN devices ON devices.id = pons.device_id
		WHERE locations.state = $1 AND measurement_ont.date BETWEEN $2 AND $3
		GROUP BY devices.sys_name, hour
		ORDER BY devices.sys_name, hour`
	err := repo.db.SelectContext(ctx, &res, query, state, initDate, endDate)
	return res, err
}

func (repo ontRepository) GetOntStatusByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error) {
	var res []entity.OntStatusCountsByState
	query := `
		SELECT 
			devices.sys_name AS sysname,
			DATE_TRUNC('hour', measurement_ont.date) AS hour,
			COUNT(DISTINCT measurement_ont.interface_id) AS pons_count,
			COUNT(CASE WHEN control_run_status = 1 THEN 1 END) AS active_count,
			COUNT(CASE WHEN control_run_status = 2 THEN 1 END) AS inactive_count,
			COUNT(CASE WHEN control_run_status NOT IN (1,2) THEN 1 END) AS unknown_count,
			COUNT(*) AS total_count
		FROM measurement_ont
		INNER JOIN fats_pon ON fats_pon.interface_id = measurement_ont.interface_id
		INNER JOIN fats ON fats.id = fats_pon.fat_id
		INNER JOIN locations ON locations.id = fats.location_id
		INNER JOIN pons ON pons.id = measurement_ont.interface_id
		INNER JOIN devices ON devices.id = pons.device_id
		WHERE locations.state = $1 AND fats.odn = $2 AND measurement_ont.date BETWEEN $3 AND $4
		GROUP BY devices.sys_name, hour
		ORDER BY devices.sys_name, hour`

	err := repo.db.SelectContext(ctx, &res, query, state, odn, initDate, endDate)
	return res, err
}

func (repo ontRepository) GetTrafficOnt(ctx context.Context, interfaceID, idx string, initDate, endDate time.Time) ([]entity.TrafficOnt, error) {
	var res []entity.TrafficOnt
	query := `
		SELECT
			hour,
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
				date AS hour,
				despt,
				serial_number,
        			line_prof_name,
				olt_distance,
        			control_mac_count, 
				control_run_status,
				bytes_in AS prev_bytes_in,
				bytes_out AS prev_bytes_out,
				LEAD(bytes_in) OVER (PARTITION BY interface_id ORDER BY date) AS curr_bytes_in,
				LEAD(bytes_out) OVER (PARTITION BY interface_id ORDER BY date) AS curr_bytes_out,
				EXTRACT(EPOCH FROM (LEAD(date) OVER (PARTITION BY interface_id ORDER BY date) - date)) AS time_diff
			FROM measurement_ont
			WHERE interface_id = $1 AND idx = $2 AND bytes_in > 0 AND bytes_out > 0 AND date BETWEEN $3 AND $4
			ORDER BY date
		) AS subquery`

	err := repo.db.SelectContext(ctx, &res, query, interfaceID, idx, initDate, endDate)
	return res, err
}
