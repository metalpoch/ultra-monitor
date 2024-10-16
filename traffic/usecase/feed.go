package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/traffic/repository"
	"gorm.io/gorm"
)

type feedUsecase struct {
	repo     repository.FeedRepository
	telegram tracking.Telegram
}

func NewFeedUsecase(db *gorm.DB, telegram tracking.Telegram) *feedUsecase {
	return &feedUsecase{repository.NewFeedRepository(db), telegram}
}

func (use feedUsecase) GetDevice(id uint) (*model.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetDevice(ctx, id)
	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetDevice - use.repo.GetDevice(ctx, %d)", id),
			err,
		)
		return nil, err
	}

	return (*model.Device)(res), err
}

func (use feedUsecase) GetAllDevice() ([]*model.DeviceLite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetAllDevice(ctx)
	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			"(feedUsecase).GetAllDevice - use.repo.GetAllDevice(ctx, %d)",
			err,
		)
		return nil, err
	}

	var devices []*model.DeviceLite
	for _, d := range res {
		devices = append(devices, &model.DeviceLite{
			ID:          d.ID,
			IP:          d.IP,
			SysName:     d.SysName,
			SysLocation: d.SysLocation,
			IsAlive:     d.IsAlive,
			LastCheck:   d.LastCheck,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
	}

	return devices, err
}

func (use feedUsecase) GetDeviceByState(state string) ([]*model.DeviceLite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetDeviceByState(ctx, state)
	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetDeviceByState - use.repo.GetDeviceByState(ctx, %s)", state),
			err,
		)
		return nil, err
	}

	var devices []*model.DeviceLite
	for _, d := range res {
		devices = append(devices, &model.DeviceLite{
			ID:          d.ID,
			IP:          d.IP,
			SysName:     d.SysName,
			SysLocation: d.SysLocation,
			IsAlive:     d.IsAlive,
			LastCheck:   d.LastCheck,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
	}

	return devices, err
}

func (use feedUsecase) GetDeviceByCounty(state, county string) ([]*model.DeviceLite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetDeviceByCounty(ctx, state, county)
	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetDeviceByCounty - use.repo.GetDeviceByCounty(ctx, %s, %s)", state, county),
			err,
		)
		return nil, err
	}

	var devices []*model.DeviceLite
	for _, d := range res {
		devices = append(devices, &model.DeviceLite{
			ID:          d.ID,
			IP:          d.IP,
			SysName:     d.SysName,
			SysLocation: d.SysLocation,
			IsAlive:     d.IsAlive,
			LastCheck:   d.LastCheck,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
	}

	return devices, err
}

func (use feedUsecase) GetDeviceByMunicipality(state, county, municipality string) ([]*model.DeviceLite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetDeviceByMunicipality(ctx, state, county, municipality)
	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetDeviceByMunicipality - use.repo.GetDeviceByMunicipality(ctx, %s, %s, %s)", state, county, municipality),
			err,
		)
		return nil, err
	}

	var devices []*model.DeviceLite
	for _, d := range res {
		devices = append(devices, &model.DeviceLite{
			ID:          d.ID,
			IP:          d.IP,
			SysName:     d.SysName,
			SysLocation: d.SysLocation,
			IsAlive:     d.IsAlive,
			LastCheck:   d.LastCheck,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
		})
	}

	return devices, err
}

func (use feedUsecase) GetInterface(id uint) (*model.Interface, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetInterface(ctx, id)
	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetInterface - use.repo.GetInterface(ctx, %d)", id),
			err,
		)
		return nil, err
	}

	return (*model.Interface)(res), err
}

func (use feedUsecase) GetInterfacesByDevice(id uint) ([]*model.InterfaceLite, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetInterfacesByDevice(ctx, id)
	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetInterface - use.repo.GetInterface(ctx, %d)", id),
			err,
		)
		return nil, err
	}

	var interfaces []*model.InterfaceLite
	for _, i := range res {
		interfaces = append(interfaces, &model.InterfaceLite{
			ID:        i.ID,
			IfIndex:   i.IfIndex,
			IfName:    i.IfName,
			IfDescr:   i.IfDescr,
			IfAlias:   i.IfAlias,
			CreatedAt: i.CreatedAt,
			UpdatedAt: i.UpdatedAt,
		})
	}

	return interfaces, err
}

func (use feedUsecase) GetLocationStates() ([]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetLocationStates(ctx)
	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			"(feedUsecase).GetLocationStates - use.repo.GetLocationStates(ctx)",
			err,
		)
		return nil, err
	}

	return res, err
}

func (use feedUsecase) GetLocationCounties(state string) ([]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetLocationCounties(ctx, state)
	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetLocationCounties - use.repo.GetLocationCounties(ctx, %s)", state),
			err,
		)
		return nil, err
	}

	return res, err
}

func (use feedUsecase) GetLocationMunicipalities(state, county string) ([]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.GetLocationMunicipalities(ctx, state, county)
	if err != nil {
		go use.telegram.Notification(
			constants.MODULE_TRAFFIC,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetLocationMunicipalities - use.repo.GetLocationMunicipalities(ctx, %s, %s)", state, county),
			err,
		)
		return nil, err
	}

	return res, err
}
