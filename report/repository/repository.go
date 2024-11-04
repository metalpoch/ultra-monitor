package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type FatRepository interface {
	Add(ctx context.Context, fat *entity.Fat) error
	Get(ctx context.Context, id uint) (*entity.Fat, error)
	Delete(ctx context.Context, id uint) error
	GetAll(ctx context.Context) ([]*entity.Fat, error)
	GetByFat(ctx context.Context, fat string) (*entity.Fat, error)
}

type ReportRepository interface {
	Add(ctx context.Context, static *entity.Report) error
	Get(ctx context.Context, id string) (*entity.Report, error)
	GetReports(ctx context.Context, category string, user_id uint) ([]*entity.Report, error)
	GetCategories(ctx context.Context) ([]*string, error)
}
