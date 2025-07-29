package usecase

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/repository"
)

type PrometheusUsecase struct {
	repo repository.PrometheusRepository
}

func NewPrometheusUsecase(db *sqlx.DB) *PrometheusUsecase {
	return &PrometheusUsecase{repository.NewPrometheusRepository(db)}
}

func (use *PrometheusUsecase) Upsert(ctx context.Context, data dto.Prometheus) error {
	return use.repo.Upsert(ctx, (entity.PrometheusUpsert)(data))

}
