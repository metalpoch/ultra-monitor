package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/measurement/model"
	"github.com/metalpoch/olt-blueprint/measurement/repository"
	"gorm.io/gorm"
)

type measurementUsecase struct {
	repo     repository.MeasurementRepository
	telegram tracking.SmartModule
}

func NewMeasurementUsecase(db *gorm.DB, telegram tracking.SmartModule) *measurementUsecase {
	return &measurementUsecase{repository.NewMeasurementRepository(db), telegram}
}

func (use measurementUsecase) Get(id uint) (*model.Measurement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	measurement := new(entity.Measurement)
	if err := use.repo.Get(ctx, id, measurement); err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(measurementUsecase).Get - use.repo.Get(ctx, %d, %v)", id, *measurement),
			err,
		)
		return &model.Measurement{}, err
	}

	return (*model.Measurement)(measurement), nil
}

func (use measurementUsecase) Upsert(measurement *model.Measurement) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := use.repo.Upsert(ctx, (*entity.Measurement)(measurement))
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(measurementUsecase).Upsert - use.repo.Upsert(ctx, %v)", *(*entity.Measurement)(measurement)),
			err,
		)
	}
	return err
}
