package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/gorm"
)

type feedRepository struct {
	db *gorm.DB
}

func NewFeedRepository(db *gorm.DB) *feedRepository {
	return &feedRepository{db}
}

func (repo feedRepository) GetDevice(ctx context.Context, id uint) (*entity.Device, error) {
	d := new(entity.Device)
	err := repo.db.WithContext(ctx).Preload("Template").First(d, id).Error
	return d, err
}

func (repo feedRepository) GetInterface(ctx context.Context, id uint) (*entity.Interface, error) {
	i := new(entity.Interface)
	err := repo.db.WithContext(ctx).Preload("Device").Preload("Device.Template").First(i, id).Error
	return i, err
}

func (repo feedRepository) GetAllDevice(ctx context.Context) ([]*entity.Device, error) {
	var devices []*entity.Device
	err := repo.db.WithContext(ctx).Find(&devices).Error
	return devices, err
}
