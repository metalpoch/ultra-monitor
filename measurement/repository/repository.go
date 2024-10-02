package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/measurement/entity"
)

type TemplateRepository interface {
	Add(ctx context.Context, template entity.Template) error
	Get(ctx context.Context, id uint) (*entity.Template, error)
	GetAll(ctx context.Context) ([]entity.Template, error)
	Update(ctx context.Context, template *entity.Template) error
	Delete(ctx context.Context, id uint) error
}

type DeviceRepository interface {
	Add(ctx context.Context, device *entity.Device) error
	Check(ctx context.Context, device *entity.Device) error
	Get(ctx context.Context, id uint) (*entity.Device, error)
	GetAll(ctx context.Context) ([]*entity.Device, error)
	GetDeviceWithOIDRows(ctx context.Context) ([]*entity.DeviceWithOID, error)
	Update(ctx context.Context, device *entity.Device) error
	Delete(ctx context.Context, id uint) error
}

type InterfaceRepository interface {
	Get(ctx context.Context, id uint) (*entity.Interface, error)
	Upsert(ctx context.Context, element *entity.Interface) error
	GetAll(ctx context.Context) ([]*entity.Interface, error)
	GetAllByDevice(ctx context.Context, id uint) ([]*entity.Interface, error)
}

type MeasurementRepository interface {
	Get(ctx context.Context, id uint, measurement *entity.Measurement) error
	Upsert(ctx context.Context, measurement *entity.Measurement) error
}

type TrafficRepository interface {
	Add(ctx context.Context, traffic *entity.Traffic) error
}
