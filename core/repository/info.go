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

func (repo infoRepository) GetDeviceBySysname(ctx context.Context, sysname string) (*entity.Device, error) {
	d := new(entity.Device)
	err := repo.db.WithContext(ctx).Preload("Template").Where("sys_name = ?", sysname).First(d).Error
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
		Joins("JOIN fat_interfaces AS fi ON interfaces.id = fi.interface_id").
		Joins("JOIN fats ON fats.id = fi.fat_id").
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
		Joins("JOIN fat_interfaces AS fi ON interfaces.id = fi.interface_id").
		Joins("JOIN fats ON fats.id = fi.fat_id").
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
		Joins("JOIN fat_interfaces AS fi ON interfaces.id = fi.interface_id").
		Joins("JOIN fats ON fats.id = fi.fat_id").
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

func (repo infoRepository) GetLocation(ctx context.Context, id uint) (*entity.Location, error) {
	d := new(entity.Location)
	err := repo.db.WithContext(ctx).First(d, id).Error
	return d, err
}

func (repo infoRepository) GetFat(ctx context.Context, id uint) (*entity.Fat, error) {
	f := new(entity.Fat)

	err := repo.db.WithContext(ctx).First(f, id).Error
	if err != nil {
		return nil, err
	}
	return f, err
}

func (repo infoRepository) GetODN(ctx context.Context, odn string) ([]*entity.FatInterface, error) {
	var fats []*entity.FatInterface

	err := repo.db.
		WithContext(ctx).
		Joins("JOIN fats ON fats.id = fat_interfaces.fat_id").
		Where("fats.odn = ?", odn).
		Find(&fats).Error
	if err != nil {
		return nil, err
	}

	return fats, nil
}

func (repo infoRepository) GetODNStates(ctx context.Context, state string) ([]*string, error) {
	var onds []*string

	err := repo.db.
		WithContext(ctx).
		Model(&entity.Fat{}).
		Joins("JOIN locations ON fats.location_id = locations.id").
		Where("locations.state = ?", state).
		Distinct().
		Pluck("fats.odn", &onds).
		Error
	if err != nil {
		return nil, err
	}

	return onds, err
}

func (repo infoRepository) GetODNStatesContries(ctx context.Context, state, country string) ([]*string, error) {
	var onds []*string

	err := repo.db.
		WithContext(ctx).
		Model(&entity.Fat{}).
		Joins("JOIN locations ON fats.location_id = locations.id").
		Where("locations.state = ? AND locations.county = ?", state, country).
		Distinct().
		Pluck("fats.odn", &onds).
		Error
	if err != nil {
		return nil, err
	}

	return onds, nil
}

func (repo infoRepository) GetODNStatesContriesMunicipality(ctx context.Context, state, country, municipality string) ([]*string, error) {
	var onds []*string

	err := repo.db.
		WithContext(ctx).
		Model(&entity.Fat{}).
		Joins("JOIN locations ON fats.location_id = locations.id").
		Where("locations.state = ? AND locations.county = ? AND locations.municipality",
			state, country, municipality).
		Distinct().
		Pluck("fats.odn", &onds).
		Error
	if err != nil {
		return nil, err
	}

	return onds, nil
}

func (repo infoRepository) GetODNDevice(ctx context.Context, id uint) ([]*string, error) {
	var odn []*string

	err := repo.db.
		WithContext(ctx).
		Model(&entity.Interface{}).
		Where("device_id = ?", id).
		Joins("INNER JOIN fat_interfaces ON interfaces.id = fat_interfaces.interface_id").
		Joins("INNER JOIN fats ON fat_interfaces.fat_id = fats.id").
		Distinct().
		Select("fats.odn").
		Scan(&odn).Error
	if err != nil {
		return nil, err
	}

	return odn, err
}

func (repo infoRepository) GetODNDevicePort(ctx context.Context, id uint, pattern string) ([]*string, error) {
	var odn []*string

	err := repo.db.WithContext(ctx).
		Model(&entity.Interface{}).
		Where("device_id = ? AND if_name LIKE ?", id, pattern).
		Joins("INNER JOIN fat_interfaces ON interfaces.id = fat_interfaces.interface_id").
		Joins("INNER JOIN fats ON fat_interfaces.fat_id = fats.id").
		Distinct().
		Pluck("fats.odn", &odn).
		Error
	if err != nil {
		return nil, err
	}
	return odn, err
}
