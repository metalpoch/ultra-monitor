package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/update/entity"
)

type TemplateRepository interface {
	Add(ctx context.Context, template *entity.Template) error
	GetByID(ctx context.Context, id uint) (*entity.Template, error)
	GetAll(ctx context.Context) ([]*entity.Template, error)
}

type DeviceRepository interface {
	Add(ctx context.Context, device *entity.Device) (uint, error)
	GetAll(ctx context.Context) ([]*entity.Device, error)
}
