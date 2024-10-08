package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/report/repository"
	"gorm.io/gorm"
)

type fatUsecase struct {
	repo     repository.FatRepository
	telegram tracking.Telegram
}

func NewFatUsecase(db *gorm.DB, telegram tracking.Telegram) *fatUsecase {
	return &fatUsecase{repository.NewFatRepository(db), telegram}
}

func (use fatUsecase) Add(fat *model.NewFat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newFat := &entity.Fat{
		Fat:         fat.Fat,
		Splitter:    fat.Splitter,
		Address:     fat.Address,
		Latitude:    fat.Latitude,
		Longitude:   fat.Longitude,
		InterfaceID: fat.InterfaceID,
	}
	err := use.repo.Add(ctx, newFat)

	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(fatUsecase).Add - use.repo.Add(ctx, %v)", newFat),
			err,
		)
	}

	return err
}

func (use fatUsecase) Get(id uint) (*model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.Get(ctx, id)
	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetDevice - use.repo.GetDevice(ctx, %d)", id),
			err,
		)
		return nil, err
	}

	return (*model.Fat)(res), nil
}
