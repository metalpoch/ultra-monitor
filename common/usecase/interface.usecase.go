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

type InterfaceUsecase struct {
	repo     repository.InterfaceRepository
	telegram tracking.SmartModule
}

func NewInterfaceUsecase(db *gorm.DB, telegram tracking.SmartModule) *InterfaceUsecase {
	return &InterfaceUsecase{repository.NewInterfaceRepository(db), telegram}
}

func (use InterfaceUsecase) Upsert(element *model.Interface) (uint64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	id, err := use.repo.Upsert(ctx, (*entity.Interface)(element))
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(interfaceUsecase).Upsert - use.repo.Upsert(ctx, %v)", *(*entity.Interface)(element)),
			err,
		)
	}

	return id, err
}

func (use InterfaceUsecase) GetAllByDevice(id uint64) ([]*model.Interface, error) {
	var interfaces []*model.Interface
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetAllByDevice(ctx, id)
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(interfaceUsecase).GetAllByDevice - use.repo.GetAllByDevice(ctx, %d)", id),
			err,
		)
		return nil, err
	}

	for _, e := range res {
		interfaces = append(interfaces, (*model.Interface)(e))
	}

	return interfaces, nil
}

func (use InterfaceUsecase) GetAllByDeviceAndIfindex(deviceID, idx uint64) (*model.Interface, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pon, err := use.repo.GetAllByDeviceAndIfindex(ctx, deviceID, idx)
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(interfaceUsecase).GetAllByDeviceAndIfindex - use.repo.GetAllByDevice(ctx, %d, %d)", deviceID, idx),
			err,
		)
		return nil, err
	}
	return (*model.Interface)(pon), nil
}

func (use InterfaceUsecase) GetAll() ([]*model.Interface, error) {
	var interfaces []*model.Interface
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetAll(ctx)
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			"(interfaceUsecase).GetAll - use.repo.GetAll(ctx)",
			err,
		)
		return nil, err
	}

	for _, e := range res {
		interfaces = append(interfaces, (*model.Interface)(e))
	}

	return interfaces, nil
}
