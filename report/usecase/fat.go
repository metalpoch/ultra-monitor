package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/osm"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	commonRepository "github.com/metalpoch/olt-blueprint/common/repository"
	"github.com/metalpoch/olt-blueprint/report/repository"
	"github.com/metalpoch/olt-blueprint/report/utils"
	"gorm.io/gorm"
)

type fatUsecase struct {
	fat      repository.FatRepository
	location commonRepository.LocationRepository
	telegram tracking.Telegram
}

func NewFatUsecase(db *gorm.DB, telegram tracking.Telegram) *fatUsecase {
	return &fatUsecase{
		fat:      repository.NewFatRepository(db),
		location: commonRepository.NewLocationRepository(db),
		telegram: telegram,
	}
}

func (use fatUsecase) Add(fat *model.NewFat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	loc, err := osm.LocationByCoord(fat.Latitude, fat.Longitude)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_OSM,
			fmt.Sprintf("(fatUsecase).Add - osm.LocationByCoord(%f, %f)", fat.Latitude, fat.Longitude),
			err,
		)
		return err
	}
	location := entity.Location{
		State:        loc.State,
		County:       loc.County,
		Municipality: loc.Municipality,
	}

	err = use.location.Find(ctx, &location)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if err := use.location.Add(ctx, &location); err != nil {
			go use.telegram.SendMessage(
				constants.MODULE_REPORT,
				constants.CATEGORY_DATABASE,
				fmt.Sprintf("(fatUsecase).Add - use.location.Add(ctx, %v)", location),
				err,
			)
			return err
		}
	} else if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(fatUsecase).Add - use.location.Find(ctx, %v)", location),
			err,
		)
		return err
	}

	newFat := &entity.Fat{
		OND:         fat.ODN,
		Fat:         fat.Fat,
		Splitter:    fat.Splitter,
		Address:     fat.Address,
		Latitude:    fat.Latitude,
		Longitude:   fat.Longitude,
		InterfaceID: fat.InterfaceID,
		LocationID:  location.ID,
	}

	err = use.fat.Add(ctx, newFat)
	if err != nil {
		go use.telegram.SendMessage(
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

	res, err := use.fat.Get(ctx, id)
	if err != nil && err != gorm.ErrRecordNotFound {
		go use.telegram.SendMessage(
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

	res, err := use.fat.GetAll(ctx)
	if err != nil && err != gorm.ErrRecordNotFound {
		go use.telegram.SendMessage(
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

	err := use.fat.Delete(ctx, id)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(fatUsecase).Delete - use.repo.Delete(ctx, %d)", id),
			err,
		)
	}

	return err
}
