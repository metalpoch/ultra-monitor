package repository

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/measurement/constants"
	"github.com/metalpoch/olt-blueprint/measurement/entity"
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
func (repo deviceRepository) GetAll(ctx context.Context) ([]*entity.Device, error) {
	var devices []*entity.Device
	err := repo.db.WithContext(ctx).Find(&devices).Error
	return devices, err
}

func (repo deviceRepository) GetDeviceWithOIDRows(ctx context.Context) ([]*entity.DeviceWithOID, error) {
	var devices []*entity.DeviceWithOID
	err := repo.db.Table(constants.TABLE_DEVICES).
		Select(constants.SELECT_TEMPLATES_ON_DEVICES).
		Joins(constants.JOIN_TEMPLATES_ON_DEVICES).
		Scan(&devices).Error

	return devices, err
}
