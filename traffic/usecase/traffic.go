package usecase

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/traffic/model"
	"github.com/metalpoch/olt-blueprint/traffic/repository"
	"gorm.io/gorm"
)

type trafficUsecase struct {
	secret []byte
	repo   repository.TrafficRepository
}

func NewTrafficUsecase(db *gorm.DB, secret []byte) *trafficUsecase {
	return &trafficUsecase{secret, repository.NewUserRepository(db)}
}

func (use trafficUsecase) GetTrafficByInterface(id uint, date *model.RangeDate) ([]*model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetTrafficByInterface(ctx, id, date)
	if err != nil {
		return nil, err
	}

	traffics := []*model.Traffic{}
	for _, e := range res {
		traffics = append(traffics, &model.Traffic{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
		})
	}

	return traffics, err
}

func (use trafficUsecase) GetTrafficByDevice(id uint, date *model.RangeDate) ([]*model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetTrafficByDevice(ctx, id, date)
	if err != nil {
		return nil, err
	}

	traffics := []*model.Traffic{}
	for _, e := range res {
		traffics = append(traffics, &model.Traffic{
			Date:      e.Date,
			Bandwidth: e.Bandwidth,
			In:        e.In,
			Out:       e.Out,
		})
	}

	return traffics, err
}
