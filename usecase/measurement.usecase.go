package usecase

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/metalpoch/ultra-monitor/repository"
)

type MeasurementUsecase struct {
	repo repository.MeasurementRepository
}

func NewMeasurementUsecase(db *sqlx.DB) *MeasurementUsecase {
	return &MeasurementUsecase{repository.NewMeasurementRepository(db)}
}

func (uc *MeasurementUsecase) UpsertOlt(olt model.Olt) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return uc.repo.UpsertOlt(ctx, (entity.Olt)(olt))
}

func (uc *MeasurementUsecase) UpsertPon(element model.Pon) (int32, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, err := uc.repo.UpsertPon(ctx, (entity.Pon)(element))
	return id, err
}

func (uc *MeasurementUsecase) GetTemportalMeasurementPon(id int32) (*model.MeasurementPon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetTemportalMeasurementPon(ctx, id)
	if err != nil {
		return nil, err
	}

	return (*model.MeasurementPon)(&res), nil
}

func (uc *MeasurementUsecase) UpsertTemportalMeasurementPon(measurement model.MeasurementPon) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return uc.repo.UpsertTemportalMeasurementPon(ctx, (entity.MeasurementPon)(measurement))
}

func (uc *MeasurementUsecase) InsertTrafficPon(measurement model.TrafficPon) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return uc.repo.InsertTrafficPon(ctx, (entity.TrafficPon)(measurement))
}

func (uc *MeasurementUsecase) InsertManyOnt(measurements []model.MeasurementOnt) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var data []entity.MeasurementOnt
	for _, measurement := range measurements {
		data = append(data, (entity.MeasurementOnt)(measurement))
	}

	return uc.repo.InsertManyMeasurementOnt(ctx, data)
}
