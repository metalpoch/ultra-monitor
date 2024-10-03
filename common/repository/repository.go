package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type DeviceRepository interface {
	Add(ctx context.Context, device *entity.Device) error
	Check(ctx context.Context, device *entity.Device) error
	Get(ctx context.Context, id uint) (*entity.Device, error)
	GetAll(ctx context.Context) ([]*entity.Device, error)
	GetDeviceWithOIDRows(ctx context.Context) ([]*entity.DeviceWithOID, error)
	Update(ctx context.Context, device *entity.Device) error
	Delete(ctx context.Context, id uint) error
}
