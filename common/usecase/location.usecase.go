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

type LocationUsecase struct {
	repo     repository.LocationRepository
	telegram tracking.SmartModule
}

func NewLocationUsecase(db *gorm.DB, telegram tracking.SmartModule) *LocationUsecase {
	return &LocationUsecase{repository.NewLocationRepository(db), telegram}
}

func (use LocationUsecase) Add(location *model.Location) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	loc := (*entity.Location)(location)
	err := use.repo.Add(ctx, loc)
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(interfaceUsecase).Upsert - use.repo.Upsert(ctx, %v)", *(*entity.Location)(location)),
			err,
		)
	}

	return loc.ID, err
}

func (use LocationUsecase) FindID(location *model.Location) (uint, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	loc := (*entity.Location)(location)
	err := use.repo.Find(ctx, loc)
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(interfaceUsecase).Upsert - use.repo.Upsert(ctx, %v)", *(*entity.Location)(location)),
			err,
		)
	}

	return loc.ID, err
}
