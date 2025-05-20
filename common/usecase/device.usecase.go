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

type DeviceUsecase struct {
	repo     repository.DeviceRepository
	telegram tracking.SmartModule
}

func NewDeviceUsecase(db *gorm.DB, telegram tracking.SmartModule) *DeviceUsecase {
	return &DeviceUsecase{repository.NewDeviceRepository(db), telegram}
}

func (use DeviceUsecase) Add(device *model.Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	newDevice := &entity.Device{
		IP:          device.IP,
		SysName:     device.SysName,
		SysLocation: device.SysLocation,
		Community:   device.Community,
		IsAlive:     device.IsAlive,
		TemplateID:  device.TemplateID,
		LastCheck:   time.Now(),
	}

	err := use.repo.Add(ctx, newDevice)
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(deviceUsecase).Add - use.repo.Add(ctx, %v)", *newDevice),
			err,
		)
	}
	return err
}

func (use DeviceUsecase) Update(id uint, device *model.AddDevice) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e, err := use.repo.GetByID(ctx, id)
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(deviceUsecase).Update - use.repo.Get(ctx, %d)", id),
			err,
		)
		return err
	}

	if device.IP != "" {
		e.IP = device.IP
	}
	if device.Community != "" {
		e.Community = device.Community
	}
	if device.Template > 0 {
		e.TemplateID = device.Template
	}

	err = use.repo.Update(ctx, e)
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(deviceUsecase).Update - use.repo.Update(ctx, %v)", e),
			err,
		)
	}

	return err
}

func (use DeviceUsecase) Delete(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := use.repo.Delete(ctx, id)
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(deviceUsecase).Add - use.repo.Delete(ctx, %v)", id),
			err,
		)
	}

	return err
}

func (use DeviceUsecase) Check(device *model.Device) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := use.repo.Check(ctx, (*entity.Device)(device))
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(deviceUsecase).Check - use.repo.Check(ctx, %v)", *(*entity.Device)(device)),
			err,
		)
	}

	return err
}

func (use DeviceUsecase) GetByID(id uint) (*model.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e, err := use.repo.GetByID(ctx, id)
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(deviceUsecase).GetByID - use.repo.GetByID(ctx, %d)\n", id),
			err,
		)
	}
	fmt.Println(e.Template)
	return (*model.Device)(e), err
}

func (use DeviceUsecase) GetAll() ([]*model.Device, error) {
	var devices []*model.Device
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetAll(ctx)
	if err != nil {
		use.telegram.SendMessage(
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

func (use DeviceUsecase) GetDeviceWithOIDRows() ([]*model.DeviceWithOID, error) {
	var devices []*model.DeviceWithOID
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	res, err := use.repo.GetDevicesWithTemplate(ctx)
	if err != nil {
		use.telegram.SendMessage(
			constants.MODULE_UPDATE,
			constants.CATEGORY_DATABASE,
			"(deviceUsecase).GetDeviceWithOIDRows - use.repo.GetDevicesWithTemplate(ctx)",
			err,
		)
		return nil, err
	}

	for _, e := range res {
		devices = append(devices, &model.DeviceWithOID{
			ID:          e.ID,
			IP:          e.IP,
			SysName:     e.SysName,
			SysLocation: e.SysLocation,
			Community:   e.Community,
			OidBw:       e.Template.OidBw,
			OidIn:       e.Template.OidIn,
			OidOut:      e.Template.OidOut,
		})
	}

	return devices, nil
}
