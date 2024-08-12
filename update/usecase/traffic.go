package usecase

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/repository"
	"github.com/metalpoch/olt-blueprint/update/utils"
)

type trafficUsecase struct {
	repo repository.TrafficRepository
}

func NewTrafficUsecase(repo repository.TrafficRepository) *trafficUsecase {
	return &trafficUsecase{repo}
}

func (use trafficUsecase) Create(count model.CountDiff) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	diffDate := count.CurrDate - count.PrevDate // la primera medicion puede ocurrir algo raro
	traffic := &entity.TrafficOLT{
		OLT:       count.OLT,
		Interface: count.Interface,
		Date:      count.CurrDate,
		KbpsIn:    utils.BytesToKbps(count.PrevBytesIn, count.CurrBytesIn, diffDate),
		KbpsOut:   utils.BytesToKbps(count.PrevBytesOut, count.CurrBytesOut, diffDate),
		Bandwidth: count.CurrBandwidth,
	}

	id, err := use.repo.Create(ctx, traffic)
	if err != nil {
		return "", err
	}

	return id, nil
}
