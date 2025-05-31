package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/constants"
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
	_, err := repo.db.NamedExecContext(ctx, constants.SQL_UPSERT_OLT, olt)
	return err
}

func (repo *measurementRepository) UpsertPon(ctx context.Context, element entity.Pon) (int32, error) {
	var id int32
	err := repo.db.QueryRowxContext(ctx, constants.SQL_UPSERT_PON, element).Scan(&id)
	return id, err
}

func (repo *measurementRepository) GetTemportalMeasurementPon(ctx context.Context, id int32) (entity.MeasurementPon, error) {
	var res entity.MeasurementPon
	err := repo.db.GetContext(ctx, &res, constants.SQL_GET_TEMPORAL_MEASUREMENT_PON, id)
	return res, err
}

func (repo *measurementRepository) UpsertTemportalMeasurementPon(ctx context.Context, measurement entity.MeasurementPon) error {
	_, err := repo.db.NamedExecContext(ctx, constants.SQL_UPSERT_TEMPORAL_MEASUREMENT_PON, measurement)
	return err
}

func (repo *measurementRepository) InsertTrafficPon(ctx context.Context, traffic entity.TrafficPon) error {
	_, err := repo.db.NamedExecContext(ctx, constants.SQL_INSERT_TRAFFIC_PON, traffic)
	return err
}

func (repo *measurementRepository) InsertManyMeasurementOnt(ctx context.Context, measurements []entity.MeasurementOnt) error {
	const fieldCount = 11
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
	query := constants.SQL_INSERT_MANY_MEASUREMENT_ONT_PREFIX + strings.Join(valueStrings, ", ")
	_, err := repo.db.ExecContext(ctx, query, valueArgs...)
	return err
}
