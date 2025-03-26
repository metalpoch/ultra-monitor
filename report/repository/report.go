package repository

import (
	"context"

	"github.com/google/uuid"
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
	rp := new(entity.Report)
	err := repo.db.WithContext(ctx).Preload("User").Where("id = ?", id).First(rp).Error
	return rp, err
}

func (repo reportRepository) GetCategories(ctx context.Context) ([]*string, error) {
	var c []*string
	err := repo.db.WithContext(ctx).Model(&entity.Report{}).Select("DISTINCT category").Pluck("category", &c).Error
	return c, err
}

func (repo reportRepository) GetReports(ctx context.Context, category string, user_id uint) ([]*entity.Report, error) {
	rp := []*entity.Report{}

	query := repo.db.WithContext(ctx).Preload("User")
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if user_id != 0 {
		query = query.Where("user_id = ?", user_id)
	}

	err := query.Find(&rp).Error

	return rp, err
}

func (repo reportRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return repo.db.WithContext(ctx).Delete(&entity.Report{}, id).Error
}
