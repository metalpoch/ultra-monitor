package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/cache"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/metalpoch/ultra-monitor/repository"
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
func (uc *TrafficUsecase) TrafficByPon(sysname string, shell, card, port int, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ifname := fmt.Sprintf("GPON %d/%d/%d", shell, card, port)
	res, err := uc.repo.TrafficByPon(ctx, sysname, ifname, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}
