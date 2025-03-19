package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/gorm"
)

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) *locationRepository {
	return &locationRepository{db}
}

func (repo locationRepository) Add(ctx context.Context, location *entity.Location) error {
	return repo.db.WithContext(ctx).Create(location).Error
}

func (repo locationRepository) Find(ctx context.Context, location *entity.Location) error {
	return repo.db.WithContext(ctx).
		Where("state = ? AND county = ? and municipality = ?", location.State, location.County, location.Municipality).
		First(location).
		Error
}
