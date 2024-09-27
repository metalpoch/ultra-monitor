package usecase

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/repository"
	"gorm.io/gorm"
)

type trafficUsecase struct {
	repo repository.TrafficRepository
}

func NewTrafficUsecase(db *gorm.DB) *trafficUsecase {
	return &trafficUsecase{repository.NewTrafficRepository(db)}
}

func (use trafficUsecase) Add(traffic *model.Traffic) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return use.repo.Add(ctx, (*entity.Traffic)(traffic))
}
