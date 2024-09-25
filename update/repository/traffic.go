package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"gorm.io/gorm"
)

type trafficRepository struct {
	db *gorm.DB
}

func NewTrafficRepository(db *gorm.DB) *trafficRepository {
	return &trafficRepository{db}
}

func (repo trafficRepository) Add(ctx context.Context, traffic *entity.Traffic) error {
	return repo.db.WithContext(ctx).Save(traffic).Error
}

func (repo trafficRepository) Get(ctx context.Context, id uint, traffic *entity.Traffic) error {
	return repo.db.WithContext(ctx).Where("interface_id = ?", id).First(traffic).Error
}
