package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/update/entity"
)

type DeviceRepository interface {
	Create(ctx context.Context, device *entity.Device) (string, error)
	FindAll(ctx context.Context) ([]*entity.Device, error)
	FindAllOffset(ctx context.Context, limit int64, offset int64) ([]*entity.Device, error)
}

type ElementRepository interface {
	FindID(ctx context.Context, olt, interfName string) (string, error)
	Create(ctx context.Context, olt *entity.ElementOLT) (string, error)
}

type CountRepository interface {
	Create(ctx context.Context, count entity.Count) (string, error)
	Find(ctx context.Context, olt, interfaceName string) (entity.Count, error)
}

type TrafficRepository interface {
	Create(ctx context.Context, traffic *entity.TrafficOLT) (string, error)
}
