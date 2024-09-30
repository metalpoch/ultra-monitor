package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/pkg/snmp"
	"github.com/metalpoch/olt-blueprint/update/repository"
	"gorm.io/gorm"
)

type deviceUsecase struct {
	repo     repository.DeviceRepository
	telegram tracking.Telegram
}

func NewDeviceUsecase(db *gorm.DB, telegram tracking.Telegram) *deviceUsecase {
	return &deviceUsecase{repository.NewDeviceRepository(db), telegram}
}

func (use deviceUsecase) Add(device *model.AddDevice) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var isAlive bool
	info, err := snmp.GetInfo(device.IP, device.Community)
	if err != nil {
		use.telegram.Notification(
			constants.MODULE_UPDATE,
			constants.CATEGORY_SNMP,
			fmt.Sprintf("(deviceUsecase).Add - snmp.GetInfo(%s, %s)", device.IP, device.Community),
			err,
		)
		isAlive = false
	} else {
		isAlive = true
	}

	newDevice := &entity.Device{
		IP:          device.IP,
		SysName:     info.SysName,
		SysLocation: info.SysLocation,
		Community:   device.Community,
		IsAlive:     isAlive,
		TemplateID:  device.Template,
		LastCheck:   time.Now(),
	}

	err = use.repo.Add(ctx, newDevice)
	if err != nil {
		use.telegram.Notification(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(deviceUsecase).Add - use.repo.Add(ctx, %v)", *newDevice),
			err,
		)
	}
	return err
}

func (use deviceUsecase) Check(device *model.Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := use.repo.Check(ctx, (*entity.Device)(device))
	if err != nil {
		use.telegram.Notification(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(deviceUsecase).Check - use.repo.Check(ctx, %v)", *(*entity.Device)(device)),
			err,
		)
	}

	return err
}

func (use deviceUsecase) GetAll() ([]*model.Device, error) {
	var devices []*model.Device
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetAll(ctx)
	if err != nil {
		use.telegram.Notification(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			"(deviceUsecase).GetAll - use.repo.GetAll(ctx)",
			err,
		)
	}

	for _, e := range res {
		devices = append(devices, (*model.Device)(e))
	}

	return devices, err
}

func (use deviceUsecase) GetDeviceWithOIDRows() ([]*model.DeviceWithOID, error) {
	var devices []*model.DeviceWithOID
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	res, err := use.repo.GetDeviceWithOIDRows(ctx)
	if err != nil {
		use.telegram.Notification(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			"(deviceUsecase).GetDeviceWithOIDRows - use.repo.GetDeviceWithOIDRows(ctx)",
			err,
		)
		return nil, err
	}

	for _, e := range res {
		devices = append(devices, (*model.DeviceWithOID)(e))
	}

	return devices, nil
}
