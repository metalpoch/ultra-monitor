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

type interfaceUsecase struct {
	repo     repository.InterfaceRepository
	telegram tracking.Telegram
}

func NewInterfaceUsecase(db *gorm.DB, telegram tracking.Telegram) *interfaceUsecase {
	return &interfaceUsecase{repository.NewInterfaceRepository(db), telegram}
}

func (use interfaceUsecase) Upsert(element *model.Interface) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := use.repo.Upsert(ctx, (*entity.Interface)(element))
	if err != nil {
		use.telegram.Notification(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(interfaceUsecase).Upsert - use.repo.Upsert(ctx, %v)", *(*entity.Interface)(element)),
			err,
		)
	}

	return err
}

func (use interfaceUsecase) GetAllByDevice(id uint) ([]*model.Interface, error) {
	var interfaces []*model.Interface
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetAllByDevice(ctx, id)
	if err != nil {
		use.telegram.Notification(
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

func (use interfaceUsecase) GetAll() ([]*model.Interface, error) {
	var interfaces []*model.Interface
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetAll(ctx)
	if err != nil {
		use.telegram.Notification(
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
