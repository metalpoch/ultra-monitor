package usecase

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/internal/cache"
	"github.com/metalpoch/olt-blueprint/internal/dto"
	"github.com/metalpoch/olt-blueprint/model"
	"github.com/metalpoch/olt-blueprint/repository"
)

type TrafficUsecase struct {
	repo  repository.TrafficRepository
	cache *cache.Redis
}

func NewTrafficUsecase(db *sqlx.DB, cache *cache.Redis) *TrafficUsecase {
	return &TrafficUsecase{repository.NewTrafficRepository(db), cache}
}

func (uc *TrafficUsecase) GetTotalTraffic(dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TotalTraffic(ctx, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *TrafficUsecase) TrafficByState(state string, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TrafficByState(ctx, state, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *TrafficUsecase) TrafficByCounty(state, county string, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TrafficByCounty(ctx, state, county, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *TrafficUsecase) TrafficByMunicipaly(state, county, municipality string, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TrafficByMunicipality(ctx, state, county, municipality, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *TrafficUsecase) TrafficByODN(state, odn string, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TrafficByODN(ctx, state, odn, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}
func (uc *TrafficUsecase) TrafficByPon(sysname, port string, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.TrafficByPon(ctx, sysname, port, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (uc *TrafficUsecase) TrafficOnt(ponID uint64, idx string, dates dto.RangeDate) ([]model.TrafficOnt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var traffic []model.TrafficOnt
	res, err := uc.repo.TrafficOnt(ctx, ponID, idx, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}
	for _, s := range res {
		traffic = append(traffic, (model.TrafficOnt)(s))
	}

	return traffic, err
}
