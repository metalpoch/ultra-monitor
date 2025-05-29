package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/entity"
)

type OntRepository interface {
	AllOntStatus(ctx context.Context, initDate, endDate time.Time) ([]entity.OntStatusCounts, error)
	GetOntStatusByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error)
	GetOntStatusByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error)
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

func (repo *ontRepository) GetOntStatusByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error) {
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

func (repo *ontRepository) GetOntStatusByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.OntStatusCountsByState, error) {
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
