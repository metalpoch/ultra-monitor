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

type MeasurementUsecase struct {
	repo     repository.MeasurementRepository
	telegram tracking.SmartModule
}

func NewMeasurementUsecase(db *gorm.DB, telegram tracking.SmartModule) *MeasurementUsecase {
	return &MeasurementUsecase{repository.NewMeasurementRepository(db), telegram}
}

func (uc MeasurementUsecase) Get(id uint) (*model.Measurement, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	measurement := new(entity.Measurement)
	if err := uc.repo.Get(ctx, id, measurement); err != nil {
		go uc.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(measurementUsecase).Get - use.repo.Get(ctx, %d, %v)", id, *measurement),
			err,
		)
		return &model.Measurement{}, err
	}

	return (*model.Measurement)(measurement), nil
}

func (uc MeasurementUsecase) Upsert(measurement *model.Measurement) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := uc.repo.Upsert(ctx, (*entity.Measurement)(measurement))
	if err != nil {
		uc.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(measurementUsecase).Upsert - use.repo.Upsert(ctx, %v)", *(*entity.Measurement)(measurement)),
			err,
		)
	}
	return err
}
