package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type TemplateRepository interface {
	Add(ctx context.Context, template entity.Template) error
	Get(ctx context.Context, id uint) (*entity.Template, error)
	GetAll(ctx context.Context) ([]entity.Template, error)
	Update(ctx context.Context, template *entity.Template) error
	Delete(ctx context.Context, id uint) error
}

type MeasurementRepository interface {
	Get(ctx context.Context, id uint, measurement *entity.Measurement) error
	Upsert(ctx context.Context, measurement *entity.Measurement) error
}

type TrafficRepository interface {
	Add(ctx context.Context, traffic *entity.Traffic) error
}
