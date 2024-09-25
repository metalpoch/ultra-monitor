package repository

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"gorm.io/gorm"
)

type deviceRepository struct {
	db *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *deviceRepository {
	return &deviceRepository{db}
}

func (repo deviceRepository) Add(ctx context.Context, device *entity.Device) error {
	return repo.db.WithContext(ctx).Create(device).Error
}

func (repo deviceRepository) Check(ctx context.Context, device *entity.CheckDevice) error {
	device.UpdatedAt = time.Now()
	return repo.db.WithContext(ctx).
		Model(&entity.Device{}).
		Where("id = ?", device.ID).
		Updates(device).Error
}

func (repo deviceRepository) GetAll(ctx context.Context) ([]*entity.Device, error) {
	var devices []*entity.Device
	err := repo.db.WithContext(ctx).Find(&devices).Error

	return devices, err
}

func (repo deviceRepository) GetDeviceWithOIDRows(ctx context.Context) ([]*entity.DeviceWithOID, error) {
	var devices []*entity.DeviceWithOID
	err := repo.db.Table("devices").
		Joins("JOIN templates ON devices.template_id = templates.id").
		Scan(&devices).Error

	return devices, err
}
