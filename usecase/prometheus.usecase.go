package usecase

import (
	"context"
	"time"

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

func (use *PrometheusUsecase) GetSnmpIndexes(ports dto.PrometheusDeviceQuery) ([]int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var query []entity.PrometheusDeviceQuery
	for i, shell := range ports.Shell {
		query = append(query, entity.PrometheusDeviceQuery{
			IP:    ports.IP,
			Shell: shell,
			Card:  ports.Card[i],
			Port:  ports.Port[i],
		})
	}

	return use.repo.GetSnmpIndexes(ctx, query)
}

func (use *PrometheusUsecase) GponPortsStatus() (*dto.PrometheusPortStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status, err := use.repo.GponPortsStatus(ctx)
	if err != nil {
		return nil, err
	}
	return (*dto.PrometheusPortStatus)(status), nil
}

func (use *PrometheusUsecase) GponPortsStatusByRegion(region string) (*dto.PrometheusPortStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status, err := use.repo.GponPortsStatusByRegion(ctx, region)
	if err != nil {
		return nil, err
	}
	return (*dto.PrometheusPortStatus)(status), nil
}

func (use *PrometheusUsecase) GponPortsStatusByState(state string) (*dto.PrometheusPortStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	status, err := use.repo.GponPortsStatusByState(ctx, state)
	if err != nil {
		return nil, err
	}
	return (*dto.PrometheusPortStatus)(status), nil
}
