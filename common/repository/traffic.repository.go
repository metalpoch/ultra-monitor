package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/gorm"
)

type TrafficRepository interface {
	AddOlt(ctx context.Context, traffic *entity.Traffic) error
	AddOnt(ctx context.Context, traffic *entity.TrafficOnt) error
}

type trafficRepository struct {
	db *gorm.DB
}

func NewTrafficRepository(db *gorm.DB) *trafficRepository {
	return &trafficRepository{db}
}

func (repo trafficRepository) AddOlt(ctx context.Context, traffic *entity.Traffic) error {
	return repo.db.WithContext(ctx).Create(traffic).Error
}

func (repo trafficRepository) AddOnt(ctx context.Context, traffic *entity.TrafficOnt) error {
	return repo.db.WithContext(ctx).Create(traffic).Error
}
