package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/gorm"
)

type fatRepository struct {
	db *gorm.DB
}

func NewFatRepository(db *gorm.DB) *fatRepository {
	return &fatRepository{db}
}

func (repo fatRepository) Add(ctx context.Context, fat *entity.Fat) error {
	return repo.db.WithContext(ctx).Create(fat).Error
}

func (repo fatRepository) Get(ctx context.Context, id uint) (*entity.Fat, error) {
	f := new(entity.Fat)
	err := repo.db.WithContext(ctx).Preload("Interface").Preload("Interface.Device").First(f, id).Error
	return f, err
}

func (repo fatRepository) GetAll(ctx context.Context) ([]*entity.Fat, error) {
	var f []*entity.Fat
	err := repo.db.WithContext(ctx).Preload("Interface").Preload("Interface.Device").Find(&f).Error
	return f, err
}

func (repo fatRepository) Delete(ctx context.Context, id uint) error {
	return repo.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Fat{}).Error
}
