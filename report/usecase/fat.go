package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/openstreetmap"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	commonRepository "github.com/metalpoch/olt-blueprint/common/repository"
	"github.com/metalpoch/olt-blueprint/report/repository"
	"github.com/metalpoch/olt-blueprint/report/utils"
	"gorm.io/gorm"
)

type fatUsecase struct {
	fat           repository.FatRepository
	location      commonRepository.LocationRepository
	interf        commonRepository.InterfaceRepository
	telegram      tracking.SmartModule
	openstreetmap openstreetmap.OSM
}

func NewFatUsecase(db *gorm.DB, telegram tracking.SmartModule, openstreetmap openstreetmap.OSM) *fatUsecase {
	return &fatUsecase{
		fat:           repository.NewFatRepository(db),
		location:      commonRepository.NewLocationRepository(db),
		interf:        commonRepository.NewInterfaceRepository(db),
		telegram:      telegram,
		openstreetmap: openstreetmap,
	}
}

func (use fatUsecase) Add(fat *model.NewFat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var (
		interfaceID uint64
		fatID       uint
	)

	// Find interfaceID
	i, err := use.interf.Get(ctx, fat.InterfaceID)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_OSM,
			fmt.Sprintf("(fatUsecase).Add - use.interf.Get(ctx, %d)", fat.InterfaceID),
			err,
		)
		return err
	}
	interfaceID = i.ID

	// Find FatID
	f, err := use.fat.GetFatByLocation(ctx, fat.Address, fat.Latitude, fat.Longitude)
	if err != nil && err != gorm.ErrRecordNotFound {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_OSM,
			fmt.Sprintf(
				"(fatUsecase).Add - use.fat.GetFatByLocation(ctx, %s, %s, %f, %f)",
				fat.Fat,
				fat.Address,
				fat.Latitude,
				fat.Longitude,
			),
			err,
		)
		return err
	} else {
		fatID = f.ID
	}

	// Create fat and Find/Create location
	if errors.Is(err, gorm.ErrRecordNotFound) {
		var locationID uint
		loc, err := use.openstreetmap.LocationByCoord(fat.Latitude, fat.Longitude)
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
			locationID = location.ID
		} else if err != nil {
			go use.telegram.SendMessage(
				constants.MODULE_REPORT,
				constants.CATEGORY_DATABASE,
				fmt.Sprintf("(fatUsecase).Add - use.location.Find(ctx, %v)", location),
				err,
			)
			return err
		} else {
			locationID = location.ID
		}

		newFat := &entity.Fat{
			ODN:        fat.ODN,
			Fat:        fat.Fat,
			Splitter:   fat.Splitter,
			Address:    fat.Address,
			Latitude:   fat.Latitude,
			Longitude:  fat.Longitude,
			LocationID: locationID,
		}

		err = use.fat.Add(ctx, newFat)
		if err != nil {
			go use.telegram.SendMessage(
				constants.MODULE_REPORT,
				constants.CATEGORY_DATABASE,
				fmt.Sprintf("(fatUsecase).Add - use.repo.Add(ctx, %v)", newFat),
				err,
			)
			return err
		}

		fatID = newFat.ID
	}

	newFatInterface := &entity.FatInterface{FatID: fatID, InterfaceID: interfaceID}
	err = use.fat.AddInterface(ctx, newFatInterface)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(fatUsecase).Add - use.repo.AddInterface(ctx, %v)", *newFatInterface),
			err,
		)
	}

	return err
}

func (use fatUsecase) Get(id uint) (*model.FatResponse, error) {
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

	return utils.FatResponse(res), nil
}

func (use fatUsecase) GetAll(page *model.Page) ([]*model.FatResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.fat.GetAll(ctx, page.Num, page.Size)
	if err != nil && err != gorm.ErrRecordNotFound {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			"(fatUsecase).GetAll - use.repo.GetAll(ctx)",
			err,
		)
		return nil, err
	}
	var fats []*model.FatResponse
	for _, e := range res {
		fats = append(fats, utils.FatResponse(e))
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
