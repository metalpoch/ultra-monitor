package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/common/repository"
	"gorm.io/gorm"
)

type trafficUsecase struct {
	repo     repository.TrafficRepository
	telegram tracking.Telegram
}

func NewTrafficUsecase(db *gorm.DB, telegram tracking.Telegram) *trafficUsecase {
	return &trafficUsecase{repository.NewTrafficRepository(db), telegram}
}

func (use trafficUsecase) Add(traffic *model.Traffic) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := use.repo.Add(ctx, (*entity.Traffic)(traffic))
	if err != nil {
		use.telegram.Notification(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(trafficUsecase).Add - use.repo.Add(ctx, %v)", *(*entity.Traffic)(traffic)),
			err,
		)
	}
	return err
}
