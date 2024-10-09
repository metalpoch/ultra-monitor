package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type FatRepository interface {
	Add(ctx context.Context, fat *entity.Fat) error
	Get(ctx context.Context, id uint) (*entity.Fat, error)
	Update(ctx context.Context, fat *entity.Fat) error
}

type ReportRepository interface {
	Add(ctx context.Context, static *entity.Report) error
	Get(ctx context.Context, id string) (*entity.Report, error)
	GetReports(ctx context.Context, category string, user_id uint) ([]*entity.Report, error)
	GetCategories(ctx context.Context) ([]*string, error)
}
