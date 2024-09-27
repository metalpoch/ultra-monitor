package usecase

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/repository"
	"gorm.io/gorm"
)

type measurementUsecase struct {
	repo repository.MeasurementRepository
}

func NewMeasurementUsecase(db *gorm.DB) *measurementUsecase {
	return &measurementUsecase{repository.NewMeasurementRepository(db)}
}

func (use measurementUsecase) Get(id uint) (*model.Measurement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	measurement := new(entity.Measurement)
	if err := use.repo.Get(ctx, id, measurement); err != nil {
		return &model.Measurement{}, err
	}

	return (*model.Measurement)(measurement), nil
}

func (use measurementUsecase) Upsert(measurement *model.Measurement) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return use.repo.Upsert(ctx, (*entity.Measurement)(measurement))
}
