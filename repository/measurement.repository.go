package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type MeasurementRepository interface {
	UpsertOlt(ctx context.Context, olt entity.Olt) error
	UpsertPon(ctx context.Context, element entity.Pon) (int32, error)
	GetTemportalMeasurementPon(ctx context.Context, id int32) (entity.MeasurementPon, error)
	UpsertTemportalMeasurementPon(ctx context.Context, measurement entity.MeasurementPon) error
	InsertTrafficPon(ctx context.Context, traffic entity.TrafficPon) error
	InsertManyMeasurementOnt(ctx context.Context, measurements []entity.MeasurementOnt) error
}

type measurementRepository struct {
	db *sqlx.DB
}

func NewMeasurementRepository(db *sqlx.DB) *measurementRepository {
	return &measurementRepository{db}
}

func (repo *measurementRepository) UpsertOlt(ctx context.Context, olt entity.Olt) error {
	query := `
    UPDATE olts SET
        sys_name = :sys_name,
        sys_location = :sys_location,
        is_alive = :is_alive,
        last_check = :last_check,
    WHERE ip = :ip`
	_, err := repo.db.NamedExecContext(ctx, query, olt)
	return err
}

func (repo *measurementRepository) UpsertPon(ctx context.Context, element entity.Pon) (int32, error) {
	var id int32
	query := `
    INSERT INTO pons (olt_ip, if_index, if_name, if_descr, if_alias)
    VALUES ($1, $2, $3, $4, $5)
    ON CONFLICT (olt_ip, if_index) DO UPDATE SET
        if_name = EXCLUDED.if_name,
        if_descr = EXCLUDED.if_descr,
        if_alias = EXCLUDED.if_alias
    RETURNING id`
	err := repo.db.QueryRowxContext(ctx, query, element.OltIP, element.IfIndex, element.IfName, element.IfDescr, element.IfAlias).Scan(&id)
	return id, err
}

func (repo *measurementRepository) GetTemportalMeasurementPon(ctx context.Context, id int32) (entity.MeasurementPon, error) {
	var res entity.MeasurementPon
	query := `SELECT * FROM measurement_pons WHERE pon_id = $1`
	err := repo.db.GetContext(ctx, &res, query, id)
	return res, err
}

func (repo *measurementRepository) UpsertTemportalMeasurementPon(ctx context.Context, measurement entity.MeasurementPon) error {
	query := `
    INSERT INTO measurement_pons (pon_id, bandwidth, bytes_in_count, bytes_out_count, date)
    VALUES (:pon_id, :bandwidth, :bytes_in_count, :bytes_out_count, :date)
    ON CONFLICT (pon_id) DO UPDATE SET
        bandwidth = EXCLUDED.bandwidth,
        bytes_in_count = EXCLUDED.bytes_in_count,
        bytes_out_count = EXCLUDED.bytes_out_count,
        date = EXCLUDED.date`
	_, err := repo.db.NamedExecContext(ctx, query, measurement)
	return err
}

func (repo *measurementRepository) InsertTrafficPon(ctx context.Context, traffic entity.TrafficPon) error {
	query := `
    INSERT INTO traffic_pons (pon_id, date, bps_in, bps_out, bandwidth_mbps_sec, bytes_in, bytes_out)
    VALUES (:pon_id, :date, :bps_in, :bps_out, :bandwidth_mbps_sec, :bytes_in, :bytes_out)`
	_, err := repo.db.NamedExecContext(ctx, query, traffic)
	return err
}

func (repo *measurementRepository) InsertManyMeasurementOnt(ctx context.Context, measurements []entity.MeasurementOnt) error {
	const fieldCount = 11
	query := `
    INSERT INTO measurement_onts (
            pon_id, idx, despt, serial_number, line_prof_name, olt_distance,
            control_mac_count, control_run_status, bytes_in_count, bytes_out_count, date
        ) VALUES `
	valueStrings := make([]string, 0, len(measurements))
	valueArgs := make([]interface{}, 0, len(measurements)*fieldCount)

	for i, m := range measurements {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)",
			i*fieldCount+1, i*fieldCount+2, i*fieldCount+3, i*fieldCount+4, i*fieldCount+5,
			i*fieldCount+6, i*fieldCount+7, i*fieldCount+8, i*fieldCount+9, i*fieldCount+10, i*fieldCount+11))
		valueArgs = append(valueArgs,
			m.PonID, m.Idx, m.Despt, m.SerialNumber, m.LineProfName, m.OltDistance,
			m.ControlMacCount, m.ControlRunStatus, m.BytesIn, m.BytesOut, m.Date)
	}
	_, err := repo.db.ExecContext(ctx, query+strings.Join(valueStrings, ", "), valueArgs...)
	return err
}
