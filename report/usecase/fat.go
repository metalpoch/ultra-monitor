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
	"github.com/metalpoch/olt-blueprint/report/utils"
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
	if err != nil && err != gorm.ErrRecordNotFound {
		go use.telegram.Notification(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(fatUsecase).Get - use.repo.Get(ctx, %d)", id),
			err,
		)
		return nil, err
	}

	return (*model.Fat)(res), nil
}

func (use fatUsecase) GetAll() ([]model.FatResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetAll(ctx)
	if err != nil && err != gorm.ErrRecordNotFound {
		go use.telegram.Notification(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			"(fatUsecase).GetAll - use.repo.GetAll(ctx)",
			err,
		)
		return nil, err
	}
	var fats []model.FatResponse
	for _, e := range res {
		fats = append(fats, utils.FatResponse((*model.Fat)(e)))
	}
	return fats, err
}

func (use fatUsecase) Delete(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := use.repo.Delete(ctx, id)
	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(fatUsecase).Delete - use.repo.Delete(ctx, %d)", id),
			err,
		)
	}

	return err
}
