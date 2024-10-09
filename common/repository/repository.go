package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type DeviceRepository interface {
	Add(ctx context.Context, device *entity.Device) error
	Check(ctx context.Context, device *entity.Device) error
	GetByID(ctx context.Context, id uint) (*entity.Device, error)
	GetAll(ctx context.Context) ([]*entity.Device, error)
	GetDevicesWithTemplate(ctx context.Context) ([]*entity.Device, error)
	Update(ctx context.Context, device *entity.Device) error
	Delete(ctx context.Context, id uint) error
}

type InterfaceRepository interface {
	Get(ctx context.Context, id uint) (*entity.Interface, error)
	Upsert(ctx context.Context, element *entity.Interface) error
	GetAll(ctx context.Context) ([]*entity.Interface, error)
	GetAllByDevice(ctx context.Context, id uint) ([]*entity.Interface, error)
}

type TrafficRepository interface {
	Add(ctx context.Context, traffic *entity.Traffic) error
}
