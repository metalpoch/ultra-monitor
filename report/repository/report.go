package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/gorm"
)

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) *reportRepository {
	return &reportRepository{db}
}

func (repo reportRepository) Add(ctx context.Context, static *entity.Report) error {
	return repo.db.WithContext(ctx).Create(static).Error
}

func (repo reportRepository) Get(ctx context.Context, id string) (*entity.Report, error) {
	s := new(entity.Report)
	err := repo.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(s).Error
	return s, err
}
