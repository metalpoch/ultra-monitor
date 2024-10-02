package repository

import (
	"context"

	commonEntity "github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/measurement/entity"
)

type TemplateRepository interface {
	Add(ctx context.Context, template commonEntity.Template) error
	Get(ctx context.Context, id uint) (*commonEntity.Template, error)
	GetAll(ctx context.Context) ([]commonEntity.Template, error)
	Update(ctx context.Context, template *commonEntity.Template) error
	Delete(ctx context.Context, id uint) error
}

type DeviceRepository interface {
	Add(ctx context.Context, device *commonEntity.Device) error
	Check(ctx context.Context, device *commonEntity.Device) error
	Get(ctx context.Context, id uint) (*commonEntity.Device, error)
	GetAll(ctx context.Context) ([]*commonEntity.Device, error)
	GetDeviceWithOIDRows(ctx context.Context) ([]*commonEntity.DeviceWithOID, error)
	Update(ctx context.Context, device *commonEntity.Device) error
	Delete(ctx context.Context, id uint) error
}

type InterfaceRepository interface {
	Get(ctx context.Context, id uint) (*commonEntity.Interface, error)
	Upsert(ctx context.Context, element *commonEntity.Interface) error
	GetAll(ctx context.Context) ([]*commonEntity.Interface, error)
	GetAllByDevice(ctx context.Context, id uint) ([]*commonEntity.Interface, error)
}

type MeasurementRepository interface {
	Get(ctx context.Context, id uint, measurement *entity.Measurement) error
	Upsert(ctx context.Context, measurement *entity.Measurement) error
}

type TrafficRepository interface {
	Add(ctx context.Context, traffic *commonEntity.Traffic) error
}
