package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/entity"
)

type OltRepository interface {
	AddOlt(ctx context.Context, olt *entity.Olt) error
	Check(ctx context.Context, olt *entity.Olt) error
	Update(ctx context.Context, olt *entity.Olt) error
	Delete(ctx context.Context, id uint64) error
	Olts(ctx context.Context, page, limit uint8) ([]entity.Olt, error)
	OltsByState(ctx context.Context, state string, page, limit uint8) ([]entity.Olt, error)
	OltsByCounty(ctx context.Context, state, county string, page, limit uint8) ([]entity.Olt, error)
	OltsByMunicipality(ctx context.Context, state, county, municipality string, page, limit uint8) ([]entity.Olt, error)

	AddTrafficOlt(ctx context.Context, traffic *entity.TrafficOlt) error
	Traffic(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
	TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
	TrafficByCounty(ctx context.Context, state, county string, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
	TrafficByMunicipality(ctx context.Context, state, county, municipality string, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
	TrafficByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
	TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
	TrafficByPON(ctx context.Context, sysname, pon string, initDate, endDate time.Time) ([]entity.TrafficOlt, error)
}

type oltRepository struct {
	db *sqlx.DB
}

func NewOltRepository(db *sqlx.DB) *oltRepository {
	return &oltRepository{db}
}

func (repo oltRepository) Add(ctx context.Context, device *entity.Olt) error {
	query := `
        INSERT INTO olt (ip, community, sys_name, sys_location, is_alive, last_check)
        VALUES (:ip, :community, :sys_name, :sys_location, :is_alive, :last_check)
    `
	_, err := repo.db.NamedExecContext(ctx, query, device)
	return err
}

func (repo oltRepository) Update(ctx context.Context, olt *entity.Olt) error {
	query := `
        UPDATE olt SET
            ip = :ip,
            community = :community,
            sys_name = :sys_name,
            sys_location = :sys_location,
            is_alive = :is_alive,
            last_check = :last_check
        WHERE id = :id
    `
	_, err := repo.db.NamedExecContext(ctx, query, olt)
	return err
}

func (repo oltRepository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM olt WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo oltRepository) Check(ctx context.Context, olt *entity.Olt) error {
	query := `
        UPDATE olt SET
            sys_name = :sys_name,
            sys_location = :sys_location,
            is_alive = :is_alive,
            last_check = :last_check,
            updated_at = :updated_at
        WHERE id = :id`

	_, err := repo.db.NamedExecContext(ctx, query, olt)
	return err
}

func (repo oltRepository) AddTrafficOlt(ctx context.Context, traffic *entity.TrafficOlt) error {
	query := `
	INSERT INTO traffic_olt (date, mbps_in, mbps_out, bandwidth_mbps_sec, mbytes_in_sec, mbytes_out_sec)
	VALUES (:date, :mbps_in, :mbps_out, :bandwidth_mbps_sec, :mbytes_in_sec, :mbytes_out_sec)`
	_, err := repo.db.NamedExecContext(ctx, query, traffic)
	return err
}

func (repo oltRepository) Olts(ctx context.Context, page, limit uint8) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	query := `SELECT * FROM olt ORDER BY sys_name LIMIT ? OFFSET ?`
	err := repo.db.SelectContext(ctx, &res, query, limit, offset)
	return res, err
}

func (repo oltRepository) OltsByState(ctx context.Context, state string, page, limit uint8) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	query := `
		SELECT id, ip, community, sys_name, sys_location, is_alive, last_check, created_at
		FROM olt
		JOIN pons ON olt.id = pons.olt_id
		JOIN fats_pon ON p.id = fats_pon.pon_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON fats.location_id = locations.id
		WHERE locations.state = ?
		LIMIT ?	OFFSET ?`
	err := repo.db.SelectContext(ctx, &res, query, state, limit, offset)
	return res, err
}

func (repo oltRepository) OltsByCounty(ctx context.Context, state, county string, page, limit uint8) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	query := `
		SELECT id, ip, community, sys_name, sys_location, is_alive, last_check, created_at
		FROM olt
		JOIN pons ON olt.id = pons.olt_id
		JOIN fats_pon ON p.id = fats_pon.pon_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON fats.location_id = locations.id
		WHERE locations.state = ? AND locations.county = ?
		LIMIT ?	OFFSET ?`
	err := repo.db.SelectContext(ctx, &res, query, state, county, limit, offset)
	return res, err
}

func (repo oltRepository) OltsByMunicipality(ctx context.Context, state, county, municipality string, page, limit uint8) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	query := `
		SELECT id, ip, community, sys_name, sys_location, is_alive, last_check, created_at
		FROM olt
		JOIN pons ON olt.id = pons.olt_id
		JOIN fats_pon ON p.id = fats_pon.pon_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON fats.location_id = locations.id
		WHERE locations.state = ? AND locations.county = ? AND locations.municipality = ?
		LIMIT ?	OFFSET ?`
	err := repo.db.SelectContext(ctx, &res, query, state, county, municipality, limit, offset)
	return res, err
}

func (repo oltRepository) Traffic(ctx context.Context, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	query := `
		SELECT
			DATE_TRUNC('minute', traffics.date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth_mbps_sec,
			SUM(bytes_in) / 1000000 AS mbytes_in_sec,
			SUM(bytes_out) / 1000000 AS mbytes_out_sec
		FROM traffic_olt
		WHERE date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', traffics.date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, initDate, endDate)
	return res, err
}

func (repo oltRepository) TrafficByState(ctx context.Context, state string, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffic_olt
		JOIN fats_pon ON fats_pon.interface_id = traffics.interface_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON locations.id = fats.location_id
		WHERE locations.state = ? AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, initDate, endDate)
	return res, err
}

func (repo oltRepository) TrafficByCounty(ctx context.Context, state, county string, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffic_olt
		JOIN fats_pon ON fats_pon.interface_id = traffics.interface_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON locations.id = fats.location_id
		WHERE locations.state = ? AND locations.county AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, county, initDate, endDate)
	return res, err
}

func (repo oltRepository) TrafficByMunicipality(ctx context.Context, state, county, municipality string, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffic_olt
		JOIN fats_pon ON fats_pon.interface_id = traffics.interface_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON locations.id = fats.location_id
		WHERE locations.state = ? AND locations.county AND locations.municipality AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, county, municipality, initDate, endDate)
	return res, err
}

func (repo oltRepository) TrafficByODN(ctx context.Context, state, odn string, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffic_olt
		JOIN fats_pon ON fats_pon.interface_id = traffics.interface_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON locations.id = fats.location_id
		WHERE locations.state = ? AND fats.odn = ? AND date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, state, odn, initDate, endDate)
	return res, err
}

func (repo oltRepository) TrafficByOLT(ctx context.Context, sysname string, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffic_olt
		JOIN pons ON pons.id = traffics.interface_id"
		JOIN devices ON devices.id = pons.device_id
		WHERE devices.sys_name = ? AND traffics.date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := sqlx.SelectContext(ctx, repo.db, &res, query, sysname, initDate, endDate)
	return res, err
}

func (repo oltRepository) TrafficByPON(ctx context.Context, sysname, pon string, initDate, endDate time.Time) ([]entity.TrafficOlt, error) {
	var res []entity.TrafficOlt
	query := `
		SELECT
			DATE_TRUNC('minute', date) AS date,
			SUM("in") / 1000000 AS mbps_in,
			SUM(out) / 1000000 AS mbps_out,
			SUM(bandwidth) / 1000000 AS bandwidth,
			SUM(bytes_in) / 1000000 AS mbytes_in,
			SUM(bytes_out) / 1000000 AS mbytes_out
		FROM traffic_olt
		JOIN pons ON pons.id = traffics.interface_id"
		JOIN devices ON devices.id = pons.device_id
		WHERE devices.sys_name = ? AND pons.if_name = ? AND traffics.date BETWEEN ? AND ?
		GROUP BY DATE_TRUNC('minute', date)
		ORDER BY date`
	err := repo.db.SelectContext(ctx, &res, query, sysname, pon, initDate, endDate)
	return res, err
}
