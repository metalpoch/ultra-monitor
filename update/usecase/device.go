package usecase

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/pkg/snmp"
	"github.com/metalpoch/olt-blueprint/update/repository"
)

type deviceUsecase struct {
	repo repository.DeviceRepository
}

func NewDeviceUsecase(repo repository.DeviceRepository) *deviceUsecase {
	return &deviceUsecase{repo}
}

func (use deviceUsecase) Create(ip, community string) error {
	device := new(entity.Device)
	sysname, err := snmp.Sysname(ip, community)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	device.Community = community
	device.IP = ip
	device.Sysname = sysname

	if _, err := use.repo.Create(ctx, device); err != nil {
		return err
	}

	return nil
}

func (use deviceUsecase) FindAll() ([]*model.Device, error) {
	devices := []*model.Device{}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := use.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, d := range res {
		devices = append(devices, &model.Device{
			IP:        d.IP,
			Community: d.Community,
			Sysname:   d.Sysname,
		})

	}

	return devices, nil
}
