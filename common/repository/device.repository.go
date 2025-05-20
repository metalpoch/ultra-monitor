package repository

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/gorm"
)

type DeviceRepository interface {
	Add(ctx context.Context, device *entity.Device) error
	Check(ctx context.Context, device *entity.Device) error
	GetByID(ctx context.Context, id uint) (*entity.Device, error)
	GetAll(ctx context.Context) ([]*entity.Device, error)
	GetDevicesWithTemplate(ctx context.Context) ([]*entity.Device, error)
	Update(ctx context.Context, device *entity.Device) error
	Delete(ctx context.Context, id uint64) error
}

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *deviceRepository {
	return &deviceRepository{db}
}

func (repo deviceRepository) Add(ctx context.Context, device *entity.Device) error {
	return repo.db.WithContext(ctx).Create(device).Error
}

func (repo deviceRepository) Update(ctx context.Context, device *entity.Device) error {
	return repo.db.WithContext(ctx).Save(device).Error
}

func (repo deviceRepository) Delete(ctx context.Context, id uint64) error {
	return repo.db.WithContext(ctx).Delete(&entity.Device{ID: id}).Error
}

func (repo deviceRepository) Check(ctx context.Context, device *entity.Device) error {
	return repo.db.WithContext(ctx).
		Model(&entity.Device{ID: device.ID}).
		Updates(map[string]interface{}{
			"SysName":     device.SysName,
			"SysLocation": device.SysLocation,
			"IsAlive":     device.IsAlive,
			"LastCheck":   device.LastCheck,
			"UpdatedAt":   time.Now(),
		}).Error
}

func (repo deviceRepository) GetByID(ctx context.Context, id uint) (*entity.Device, error) {
	device := new(entity.Device)
	err := repo.db.WithContext(ctx).First(device, id).Error
	return device, err
}

func (repo deviceRepository) GetAll(ctx context.Context) ([]*entity.Device, error) {
	var devices []*entity.Device
	err := repo.db.WithContext(ctx).Find(&devices).Error
	return devices, err
}

func (repo deviceRepository) GetDevicesWithTemplate(ctx context.Context) ([]*entity.Device, error) {
	var devices []*entity.Device
	err := repo.db.WithContext(ctx).Preload("Template").Find(&devices).Error
	return devices, err
}
