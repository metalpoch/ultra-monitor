package usecase

import (
	"context"
	"log"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/pkg/snmp"
	"github.com/metalpoch/olt-blueprint/update/repository"
	"gorm.io/gorm"
)

type deviceUsecase struct {
	repo repository.DeviceRepository
}

func NewDeviceUsecase(db *gorm.DB) *deviceUsecase {
	return &deviceUsecase{repository.NewDeviceRepository(db)}
}

func (use deviceUsecase) Add(device *model.AddDevice) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	var isAlive bool
	info, err := snmp.GetInfo(device.IP, device.Community)
	if err != nil {
		log.Println("snmp error on try get the sysname:", err)
		isAlive = false
	} else {
		isAlive = true
	}

	if err := use.repo.Add(ctx, &entity.Device{
		IP:          device.IP,
		SysName:     info.SysName,
		SysLocation: info.SysLocation,
		Community:   device.Community,
		IsAlive:     isAlive,
		TemplateID:  device.Template,
		LastCheck:   time.Now(),
	}); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (use deviceUsecase) Check(device *model.CheckDevice) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return use.repo.Check(ctx, (*entity.CheckDevice)(device))
}

func (use deviceUsecase) GetAll() ([]*model.Device, error) {
	var devices []*model.Device
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetAll(ctx)

	// Gestionar errores (con Axios por ejemplo)
	// ...

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
		return nil, err
	}

	for _, e := range res {
		devices = append(devices, (*model.DeviceWithOID)(e))
	}

	return devices, nil
}
