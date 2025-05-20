package usecase

import (
	"context"
	"time"

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

func (uc MeasurementUsecase) GetOlt(id uint64) (*model.MeasurementOlt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	measurement := new(entity.Measurement)
	if err := uc.repo.GetOlt(ctx, id, measurement); err != nil {
		return &model.MeasurementOlt{}, err
	}

	return (*model.MeasurementOlt)(measurement), nil
}

func (uc MeasurementUsecase) UpsertOlt(measurement *model.MeasurementOlt) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return uc.repo.UpsertOlt(ctx, (*entity.Measurement)(measurement))
}

func (uc MeasurementUsecase) InsertManyOnt(measurements []model.MeasurementOnt) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data []entity.MeasurementOnt
	for _, measurement := range measurements {
		/*data = append(data, (entity.MeasurementOnt)(measurement))*/
		data = append(data, entity.MeasurementOnt{
			InterfaceID:      measurement.InterfaceID,
			Idx:              measurement.Idx,
			Despt:            measurement.Despt,
			SerialNumber:     measurement.SerialNumber,
			LineProfName:     measurement.LineProfName,
			OltDistance:      measurement.OltDistance,
			ControlMacCount:  measurement.ControlMacCount,
			ControlRunStatus: measurement.ControlRunStatus,
			BytesIn:          measurement.BytesIn,
			BytesOut:         measurement.BytesOut,
			Date:             measurement.Date,
		})
	}

	return uc.repo.InsertManyOnt(ctx, &data)
}
