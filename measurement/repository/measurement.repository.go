package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MeasurementRepository interface {
	GetOlt(ctx context.Context, id uint64, measurement *entity.Measurement) error
	UpsertOlt(ctx context.Context, measurement *entity.Measurement) error
	InsertManyOnt(ctx context.Context, measurements *[]entity.MeasurementOnt) error
}

type measurementRepository struct {
	db *gorm.DB
}

func NewMeasurementRepository(db *gorm.DB) *measurementRepository {
	return &measurementRepository{db}
}

func (repo measurementRepository) GetOlt(ctx context.Context, id uint64, measurement *entity.Measurement) error {
	return repo.db.WithContext(ctx).Where("interface_id = ?", id).First(measurement).Error
}

func (repo measurementRepository) UpsertOlt(ctx context.Context, measurement *entity.Measurement) error {
	return repo.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: constants.MEASUREMENT_OLT_COLUMN_INTERFACE_ID}},
		UpdateAll: true,
	}).Create(measurement).Error
}

func (repo measurementRepository) InsertManyOnt(ctx context.Context, measurements *[]entity.MeasurementOnt) error {
	return repo.db.WithContext(ctx).Create(measurements).Error
}
