package usecase

import (
	"context"
	"fmt"
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

func (use measurementUsecase) Get(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	measurement := new(entity.Measurement)
	fmt.Println(measurement)
	// return use.repo.Get(ctx, id, measurement)

	err := use.repo.Get(ctx, id, measurement)
	if err != nil {
		return err
	}

	fmt.Println(measurement)
	return nil
}

func (use measurementUsecase) Upsert(measurement *model.Measurement) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return use.repo.Upsert(ctx, (*entity.Measurement)(measurement))
}
