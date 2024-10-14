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

func (repo feedRepository) GetAllDevice(ctx context.Context) ([]*entity.Device, error) {
	var devices []*entity.Device
	err := repo.db.WithContext(ctx).Find(&devices).Error
	return devices, err
}

func (repo feedRepository) GetInterface(ctx context.Context, id uint) (*entity.Interface, error) {
	i := new(entity.Interface)
	err := repo.db.WithContext(ctx).Preload("Device").Preload("Device.Template").First(i, id).Error
	return i, err
}

func (repo feedRepository) GetInterfacesByDevice(ctx context.Context, id uint) ([]*entity.Interface, error) {
	var ifaces []*entity.Interface
	err := repo.db.WithContext(ctx).Where("device_id = ?", id).Find(&ifaces).Error
	return ifaces, err
}

func (repo feedRepository) GetLocationStates(ctx context.Context) ([]*string, error) {
	var l []*string
	err := repo.db.WithContext(ctx).Model(&entity.Location{}).Select("DISTINCT state").Pluck("state", &l).Error
	return l, err
}

func (repo feedRepository) GetLocationCounties(ctx context.Context, state string) ([]*string, error) {
	var l []*string
	err := repo.db.WithContext(ctx).Model(&entity.Location{}).
		Select("DISTINCT county").
		Where("state = ?", state).
		Pluck("county", &l).
		Error
	return l, err
}

func (repo feedRepository) GetLocationMunicipalities(ctx context.Context, state, county string) ([]*string, error) {
	var l []*string
	err := repo.db.WithContext(ctx).Model(&entity.Location{}).
		Select("DISTINCT municipality").
		Where("state = ? AND county = ?", state, county).
		Pluck("municipality", &l).
		Error
	return l, err
}
