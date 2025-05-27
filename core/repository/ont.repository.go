package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
	commonModel "github.com/metalpoch/olt-blueprint/common/model"
	"gorm.io/gorm"
)

type OntRepository interface {
	GetOntStatus(ctx context.Context, date *commonModel.TrafficRangeDate) ([]entity.UserStatusCounts, error)
	GetOntStatusByState(ctx context.Context, state string, date *commonModel.TrafficRangeDate) ([]entity.OntStatusCountsByState, error)
	GetOntStatusByODN(ctx context.Context, state, odn string, date *commonModel.TrafficRangeDate) ([]entity.OntStatusCountsByState, error)
	GetTrafficOnt(ctx context.Context, interfaceID, idx string, date *commonModel.TrafficRangeDate) ([]entity.OntTraffic, error)
}

type ontRepository struct {
	db *gorm.DB
}

func NewOntRepository(db *gorm.DB) *ontRepository {
	return &ontRepository{db}
}

func (repo ontRepository) GetOntStatus(ctx context.Context, date *commonModel.TrafficRangeDate) ([]entity.UserStatusCounts, error) {
	var res []entity.UserStatusCounts
	err := repo.db.WithContext(ctx).Raw(`
    		SELECT
			l.state,
        		DATE_TRUNC('hour', measurement_onts.date) AS hour,
			COUNT(DISTINCT measurement_onts.interface_id) AS pons_count,
			COUNT(CASE WHEN control_run_status = 1 THEN 1 END) AS active_count,
			COUNT(CASE WHEN control_run_status = 2 THEN 1 END) AS inactive_count,
			COUNT(CASE WHEN control_run_status NOT IN (1,2) THEN 1 END) AS unknown_count,
			COUNT(*) AS total_count
		FROM measurement_onts
		INNER JOIN fats_pon ON fats_pon.interface_id = measurement_onts.interface_id
		INNER JOIN fats ON fats.id = fats_pon.fat_id
		INNER JOIN locations AS l ON l.id = fats.location_id
		WHERE measurement_onts.date >= ? AND measurement_onts.date < ?
		GROUP BY l.state, hour
		ORDER BY l.state, hour`,
		date.InitDate,
		date.EndDate,
	).Scan(&res).Error

	return res, err
}

func (repo ontRepository) GetOntStatusByState(ctx context.Context, state string, date *commonModel.TrafficRangeDate) ([]entity.OntStatusCountsByState, error) {
	var res []entity.OntStatusCountsByState
	err := repo.db.WithContext(ctx).Raw(`
		SELECT devices.sys_name AS sysname,
			DATE_TRUNC('hour', measurement_onts.date) AS hour,
			COUNT(DISTINCT measurement_onts.interface_id) AS pons_count,
			COUNT(CASE WHEN control_run_status = 1 THEN 1 END) AS active_count,
			COUNT(CASE WHEN control_run_status = 2 THEN 1 END) AS inactive_count,
			COUNT(CASE WHEN control_run_status NOT IN (1,2) THEN 1 END) AS unknown_count,
			COUNT(*) AS total_count
		FROM measurement_onts
		INNER JOIN fats_pon ON fats_pon.interface_id = measurement_onts.interface_id
		INNER JOIN fats ON fats.id = fats_pon.fat_id
		INNER JOIN locations ON locations.id = fats.location_id
		INNER JOIN pons ON pons.id = measurement_onts.interface_id
		INNER JOIN devices ON devices.id = pons.device_id
		WHERE locations.state = ? AND measurement_onts.date BETWEEN ? AND ?
		GROUP BY devices.sys_name, hour
		ORDER BY devices.sys_name, hour`,
		state,
		date.InitDate,
		date.EndDate,
	).Scan(&res).Error

	return res, err
}

func (repo ontRepository) GetOntStatusByODN(ctx context.Context, state, odn string, date *commonModel.TrafficRangeDate) ([]entity.OntStatusCountsByState, error) {
	var res []entity.OntStatusCountsByState
	err := repo.db.WithContext(ctx).Raw(`
		SELECT  devices.sys_name AS sysname,
			DATE_TRUNC('hour', measurement_onts.date) AS hour,
			COUNT(DISTINCT measurement_onts.interface_id) AS pons_count,
			COUNT(CASE WHEN control_run_status = 1 THEN 1 END) AS active_count,
			COUNT(CASE WHEN control_run_status = 2 THEN 1 END) AS inactive_count,
			COUNT(CASE WHEN control_run_status NOT IN (1,2) THEN 1 END) AS unknown_count,
			COUNT(*) AS total_count
		FROM measurement_onts
		INNER JOIN fats_pon ON fats_pon.interface_id = measurement_onts.interface_id
		INNER JOIN fats ON fats.id = fats_pon.fat_id
		INNER JOIN locations ON locations.id = fats.location_id
		INNER JOIN pons ON pons.id = measurement_onts.interface_id
		INNER JOIN devices ON devices.id = pons.device_id
		WHERE locations.state = ? AND fats.odn = ? AND measurement_onts.date BETWEEN ? AND ?
		GROUP BY devices.sys_name, hour
		ORDER BY devices.sys_name, hour`,
		state,
		odn,
		date.InitDate,
		date.EndDate,
	).Scan(&res).Error

	return res, err
}

func (repo ontRepository) GetTrafficOnt(ctx context.Context, interfaceID, idx string, date *commonModel.TrafficRangeDate) ([]entity.OntTraffic, error) {
	var res []entity.OntTraffic
	err := repo.db.WithContext(ctx).Raw(`
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
			FROM measurement_onts
			WHERE interface_id = ? AND idx = ? AND bytes_in > 0 AND bytes_out > 0 AND date BETWEEN ? AND ?
			ORDER BY date
		) AS subquery`,
		interfaceID,
		idx,
		date.InitDate,
		date.EndDate,
	).Scan(&res).Error

	return res, err
}
