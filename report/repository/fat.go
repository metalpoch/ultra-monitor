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

func (repo fatRepository) AddInterface(ctx context.Context, factIntf *entity.FatInterface) error {
	return repo.db.WithContext(ctx).Create(factIntf).Error
}

func (repo fatRepository) Get(ctx context.Context, id uint) (*entity.FatInterface, error) {
	f := new(entity.FatInterface)
	err := repo.db.WithContext(ctx).
		Preload("Fat").
		Preload("Fat.Location").
		Preload("Interface").
		Preload("Interface.Device").
		First(f, "fat_id = ?", id).
		Error
	return f, err
}

func (repo fatRepository) GetFatByLocation(ctx context.Context, address string, lat, lon float64) (*entity.Fat, error) {
	f := new(entity.Fat)
	err := repo.db.WithContext(ctx).
		Where("address = ? AND latitude = ? AND longitude = ?", address, lat, lon).
		First(&f).
		Error
	return f, err
}

func (repo fatRepository) GetAll(ctx context.Context, page, pageSize int) ([]*entity.FatInterface, error) {
	var f []*entity.FatInterface
	err := repo.db.WithContext(ctx).
		Preload("Fat").
		Preload("Fat.Location").
		Preload("Interface").
		Preload("Interface.Device").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Order("id asc").
		Find(&f).Error
	return f, err
}

func (repo fatRepository) Delete(ctx context.Context, id uint) error {
	return repo.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.Fat{}).Error
}

/*func (repo fatRepository) GetByFat(ctx context.Context, fat string) ([]*entity.FatInterface, error) {
 *        var f []*entity.FatInterface
 *        err := repo.db.WithContext(ctx).
 *                Preload("Fat").
 *                Preload("Fat.Location").
 *                Preload("Interface").
 *                Preload("Interface.Device").
 *                Where("fat.fat = ?", fat).
 *                First(f).
 *                Error
 *        return f, err
 *}*/
