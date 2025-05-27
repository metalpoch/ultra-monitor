package repository

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/gorm"
)

type OltRepository interface {
	Olts(ctx context.Context, page, limit uint8) ([]entity.Olt, error)
	OltsByState(ctx context.Context, state string, page, limit uint8) ([]entity.Olt, error)
	OltsByCounty(ctx context.Context, state, county string, page, limit uint8) ([]entity.Olt, error)
	OltsByMunicipality(ctx context.Context, state, county, municipality string, page, limit uint8) ([]entity.Olt, error)

	Traffic(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
	TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
	TrafficByCounty(ctx context.Context, state, county string, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
	TrafficByMunicipality(ctx context.Context, state, county, municipality string, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
	TrafficByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
	TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
	TrafficByPON(ctx context.Context, sysname, pon string, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
}

type oltRepository struct {
	db *gorm.DB
}

func NewOltRepository(db *gorm.DB) *oltRepository {
	return &oltRepository{db}
}

func (repo oltRepository) Olts(ctx context.Context, page, limit uint8) ([]entity.Olt, error) {
	offset := (page - 1) * limit
	var res []entity.Olt
	err := repo.db.WithContext(ctx).Raw(
		`SELECT * FROM olts ORDER BY sys_name LIMIT ? OFFSET ?`,
		limit,
		offset,
	).Find(&res).Error
	return res, err
}

func (repo oltRepository) OltsByState(ctx context.Context, state string, page, limit uint8) ([]entity.Olt, error) {
	offset := (page - 1) * limit
	var res []entity.Olt
	err := repo.db.WithContext(ctx).Raw(`
		SELECT 
			o.id,
			o.ip, 
			o.sys_name, 
			o.sys_location, 
			o.is_alive, 
			o.last_check, 
			o.created_at, 
		FROM olts as o
		JOIN pons AS p ON o.id = p.olt_id
		JOIN fats_pon AS fp ON p.id = fp.pon_id
		JOIN fats ON fats.id = fp.fat_id
		JOIN locations AS l ON fats.location_id = l.id
		WHERE l.state = ?
		LIMIT ?
		OFFSET ?`,
		state,
		limit,
		offset,
	).Find(&res).Error
	return res, err
}

func (repo oltRepository) OltsByCounty(ctx context.Context, state, county string, page, limit uint8) ([]entity.Olt, error) {
	offset := (page - 1) * limit
	var res []entity.Olt
	err := repo.db.WithContext(ctx).Raw(`
		SELECT 
			o.id,
			o.ip, 
			o.sys_name, 
			o.sys_location, 
			o.is_alive, 
			o.last_check, 
			o.created_at, 
		FROM olts as o
		JOIN pons AS p ON o.id = p.olt_id
		JOIN fats_pon AS fp ON p.id = fp.pon_id
		JOIN fats ON fats.id = fp.fat_id
		JOIN locations AS l ON fats.location_id = l.id
		WHERE l.state = ?  AND l.county = ?
		LIMIT ?
		OFFSET ?`,
		state,
		county,
		limit,
		offset,
	).Find(&res).Error
	return res, err
}

func (repo oltRepository) OltsByMunicipality(ctx context.Context, state, county, municipality string, page, limit uint8) ([]entity.Olt, error) {
	offset := (page - 1) * limit
	var res []entity.Olt
	err := repo.db.WithContext(ctx).Raw(`
		SELECT 
			o.id,
			o.ip, 
			o.sys_name, 
			o.sys_location, 
			o.is_alive, 
			o.last_check, 
			o.created_at, 
		FROM olts as o
		JOIN pons AS p ON o.id = p.olt_id
		JOIN fats_pon AS fp ON p.id = fp.pon_id
		JOIN fats ON fats.id = fp.fat_id
		JOIN locations AS l ON fats.location_id = l.id
		WHERE l.state = ?  AND l.county = ? AND l.municipality = ?
		LIMIT ?
		OFFSET ?`,
		state,
		county,
		municipality,
		limit,
		offset,
	).Find(&res).Error
	return res, err
}

func (repo oltRepository) Traffic(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	err := repo.db.WithContext(ctx).Raw(`
		SELECT
			DATE_TRUNC('minute', traffics.date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffics_olt
		WHERE date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', traffics.date)
		ORDER BY date`,
		initDate,
		endDate,
	).Find(&res).Error
	return res, err
}

func (repo oltRepository) TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	err := repo.db.WithContext(ctx).Raw(`
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffics_olt
		JOIN fats_pon ON fats_pon.interface_id = traffics.interface_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations AS l ON l.id = fats.location_id
		WHERE l.state = state AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`,
		state,
		initDate,
		endDate,
	).Find(&res).Error
	return res, err
}

func (repo oltRepository) TrafficByCounty(ctx context.Context, state, county string, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	err := repo.db.WithContext(ctx).Raw(`
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffics_olt
		JOIN fats_pon ON fats_pon.interface_id = traffics.interface_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations AS l ON l.id = fats.location_id
		WHERE l.state ? AND l.county = ? AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`,
		state,
		county,
		initDate,
		endDate,
	).Find(&res).Error
	return res, err
}

func (repo oltRepository) TrafficByMunicipality(ctx context.Context, state, county, municipality string, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	err := repo.db.WithContext(ctx).Raw(`
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffics_olt
		JOIN fats_pon ON fats_pon.interface_id = traffics.interface_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations AS l ON l.id = fats.location_id
		WHERE l.state ? AND l.county = ? AND l.municipality AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`,
		state,
		county,
		initDate,
		endDate,
	).Find(&res).Error
	return res, err
}

func (repo oltRepository) TrafficByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	err := repo.db.WithContext(ctx).Raw(`
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffics_olt
		JOIN fats_pon ON fats_pon.interface_id = traffics.interface_id
		JOIN fats ON fats.id = fats_pon.fat_id
		WHERE fats.odn = ? AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`,
		state,
		initDate,
		endDate,
	).Find(&res).Error
	return res, err
}

func (repo oltRepository) TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	err := repo.db.WithContext(ctx).Raw(`
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffics_olt
		JOIN pons ON pons.id = traffics.interface_id"
		JOIN devices ON devices.id = pons.device_id
		WHERE devices.sys_name = ? AND traffics.date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`,
		sysname,
		initDate,
		endDate,
	).Find(&res).Error
	return res, err
}

func (repo oltRepository) TrafficByPON(ctx context.Context, sysname, pon string, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	err := repo.db.WithContext(ctx).Raw(`
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffics_olt
		JOIN pons ON pons.id = traffics.interface_id"
		JOIN devices ON devices.id = pons.device_id
		WHERE devices.sys_name = ? AND pons.if_name = ? AND traffics.date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`,
		sysname,
		pon,
		initDate,
		endDate,
	).Find(&res).Error
	return res, err
}
