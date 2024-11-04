package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/gorm"
)

type infoRepository struct {
	db *gorm.DB
}

func NewInfoRepository(db *gorm.DB) *infoRepository {
	return &infoRepository{db}
}

func (repo infoRepository) GetDevice(ctx context.Context, id uint) (*entity.Device, error) {
	d := new(entity.Device)
	err := repo.db.WithContext(ctx).Preload("Template").First(d, id).Error
	return d, err
}

func (repo infoRepository) GetDeviceByIP(ctx context.Context, ip string) (*entity.Device, error) {
	d := new(entity.Device)
	err := repo.db.WithContext(ctx).Preload("Template").Where("ip = ?", ip).First(d).Error
	return d, err
}

func (repo infoRepository) GetAllDevice(ctx context.Context) ([]*entity.Device, error) {
	var devices []*entity.Device
	err := repo.db.WithContext(ctx).Find(&devices).Error
	return devices, err
}

func (repo infoRepository) GetDeviceByState(ctx context.Context, state string) ([]*entity.Device, error) {
	var devices []*entity.Device
	err := repo.db.WithContext(ctx).
		Joins("JOIN interfaces ON devices.id = interfaces.device_id").
		Joins("JOIN fats ON interfaces.id = fats.interface_id").
		Joins("JOIN locations ON fats.location_id = locations.id").
		Select("DISTINCT devices.id, devices.ip, devices.sys_name", "devices.sys_location", "devices.is_alive", "devices.last_check", "devices.created_at", "devices.updated_at").
		Where("locations.state = ?", state).
		Find(&devices).Error

	return devices, err
}

func (repo infoRepository) GetDeviceByCounty(ctx context.Context, state, county string) ([]*entity.Device, error) {
	var devices []*entity.Device
	err := repo.db.WithContext(ctx).
		Joins("JOIN interfaces ON devices.id = interfaces.device_id").
		Joins("JOIN fats ON interfaces.id = fats.interface_id").
		Joins("JOIN locations ON fats.location_id = locations.id").
		Select("DISTINCT devices.id, devices.ip, devices.sys_name", "devices.sys_location", "devices.is_alive", "devices.last_check", "devices.created_at", "devices.updated_at").
		Where("locations.state = ? AND locations.county = ?", state, county).
		Find(&devices).Error
	return devices, err
}

func (repo infoRepository) GetDeviceByMunicipality(ctx context.Context, state, county, municipality string) ([]*entity.Device, error) {
	var devices []*entity.Device
	err := repo.db.WithContext(ctx).
		Joins("JOIN interfaces ON devices.id = interfaces.device_id").
		Joins("JOIN fats ON interfaces.id = fats.interface_id").
		Joins("JOIN locations ON fats.location_id = locations.id").
		Select("DISTINCT devices.id, devices.ip, devices.sys_name", "devices.sys_location", "devices.is_alive", "devices.last_check", "devices.created_at", "devices.updated_at").
		Where("locations.state = ? AND locations.county = ? AND locations.municipality = ?", state, county, municipality).
		Find(&devices).Error
	return devices, err
}

func (repo infoRepository) GetInterface(ctx context.Context, id uint) (*entity.Interface, error) {
	i := new(entity.Interface)
	err := repo.db.WithContext(ctx).Preload("Device").Preload("Device.Template").First(i, id).Error
	return i, err
}

func (repo infoRepository) GetInterfacesByDevice(ctx context.Context, id uint) ([]*entity.Interface, error) {
	var ifaces []*entity.Interface
	err := repo.db.WithContext(ctx).Where("device_id = ?", id).Find(&ifaces).Error
	return ifaces, err
}

func (repo infoRepository) GetInterfacesByDeviceAndPorts(ctx context.Context, id uint, pattern string) ([]*entity.Interface, error) {
	var ifaces []*entity.Interface
	err := repo.db.WithContext(ctx).Where("device_id = ? AND if_name LIKE ?", id, pattern).Find(&ifaces).Error
	return ifaces, err
}

func (repo infoRepository) GetLocationStates(ctx context.Context) ([]*string, error) {
	var l []*string
	err := repo.db.WithContext(ctx).Model(&entity.Location{}).Select("DISTINCT state").Pluck("state", &l).Error
	return l, err
}

func (repo infoRepository) GetLocationCounties(ctx context.Context, state string) ([]*string, error) {
	var l []*string
	err := repo.db.WithContext(ctx).Model(&entity.Location{}).
		Select("DISTINCT county").
		Where("state = ?", state).
		Pluck("county", &l).
		Error
	return l, err
}

func (repo infoRepository) GetLocationMunicipalities(ctx context.Context, state, county string) ([]*string, error) {
	var l []*string
	err := repo.db.WithContext(ctx).Model(&entity.Location{}).
		Select("DISTINCT municipality").
		Where("state = ? AND county = ?", state, county).
		Pluck("municipality", &l).
		Error
	return l, err
}
