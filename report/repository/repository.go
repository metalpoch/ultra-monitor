package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/metalpoch/olt-blueprint/common/entity"
)

type FatRepository interface {
	Add(ctx context.Context, fat *entity.Fat) error
	AddInterface(ctx context.Context, factIntf *entity.FatInterface) error
	Get(ctx context.Context, id uint) (*entity.FatInterface, error)
	Delete(ctx context.Context, id uint) error
	GetAll(ctx context.Context, page, pageSize int) ([]*entity.FatInterface, error)
	/*GetByFat(ctx context.Context, fat string) ([]*entity.FatInterface, error)*/
	GetFatByLocation(ctx context.Context, address string, lat, lon float64) (*entity.Fat, error)
}

type ReportRepository interface {
	Add(ctx context.Context, static *entity.Report) error
	Get(ctx context.Context, id string) (*entity.Report, error)
	GetReports(ctx context.Context, category string, user_id uint) ([]*entity.Report, error)
	GetCategories(ctx context.Context) ([]*string, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
