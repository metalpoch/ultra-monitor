package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
	"gorm.io/gorm"
)

type TrafficRepository interface {
	GetTrafficByInterface(ctx context.Context, id uint, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByDevice(ctx context.Context, id uint, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByFat(ctx context.Context, id uint, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByLocationID(ctx context.Context, id uint, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByState(ctx context.Context, state string, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByCounty(ctx context.Context, state, county string, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByMunicipality(ctx context.Context, state, county, municipality string, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByODN(ctx context.Context, odn string, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	// TODO: replace models by entities response
	GetTotalTrafficByState(ctx context.Context, ids []uint, month string) (*model.TrafficState, error)
	GetTotalTrafficByODN(ctx context.Context, ids []uint, month string) (*model.TrafficODN, error)
}

type trafficRepository struct {
	db *gorm.DB
}

func NewTrafficRepository(db *gorm.DB) *trafficRepository {
	return &trafficRepository{db}
}

func (repo trafficRepository) GetTrafficByInterface(ctx context.Context, id uint, date *model.TranficRangeDate) ([]*entity.Traffic, error) {
	var traffic []*entity.Traffic
	err := repo.db.WithContext(ctx).
		Where("interface_id = ? AND date BETWEEN ? AND ?", id, date.InitDate, date.EndDate).
		Order("date").
		Find(&traffic).
		Error
	return traffic, err
}

func (repo trafficRepository) GetTrafficByDevice(ctx context.Context, id uint, date *model.TranficRangeDate) ([]*entity.Traffic, error) {
	var traffic []*entity.Traffic
	err := repo.db.WithContext(ctx).
		Model(&entity.Traffic{}).
		Select("date, SUM(\"in\") AS \"in\", SUM(out) AS out, SUM(bandwidth) as bandwidth, SUM(bytes_in) AS bytes_in, SUM(bytes_out) AS bytes_out").
		Joins("JOIN interfaces ON interfaces.id = traffics.interface_id").
		Where("interfaces.device_id = ? AND traffics.date BETWEEN ? AND ?", id, date.InitDate, date.EndDate).
		Group("date").
		Order("date").
		Find(&traffic).
		Error
	return traffic, err
}

func (repo trafficRepository) GetTrafficByFat(ctx context.Context, id uint, date *model.TranficRangeDate) ([]*entity.Traffic, error) {
	var traffic []*entity.Traffic
	err := repo.db.WithContext(ctx).
		Select("date, \"in\", out, bandwidth, bytes_in, bytes_out").
		Where("date BETWEEN ? AND ?", date.InitDate, date.EndDate).
		Joins("JOIN fats ON fats.interface_id = traffics.interface_id").
		Where("fats.id = ?", id).
		Order("date").
		Find(&traffic).
		Error
	return traffic, err
}

func (repo trafficRepository) GetTrafficByLocationID(ctx context.Context, id uint, date *model.TranficRangeDate) ([]*entity.Traffic, error) {
	var traffic []*entity.Traffic
	err := repo.db.WithContext(ctx).
		Model(&entity.Traffic{}).
		Select("date, SUM(\"in\") AS \"in\", SUM(out) AS out, SUM(bandwidth) as bandwidth, SUM(bytes_in) AS bytes_in, SUM(bytes_out) AS bytes_out").
		Joins("JOIN fat_interfaces ON fat_interfaces.interface_id = traffics.interface_id").
		Joins("JOIN fats ON fats.id = fat_interfaces.fat_id").
		Where("f.location_id = ? AND traffics.date BETWEEN ? AND ?", id, date.InitDate, date.EndDate).
		Group("date").
		Order("date").
		Find(&traffic).
		Error

	return traffic, err
}

func (repo trafficRepository) GetTrafficByState(ctx context.Context, state string, date *model.TranficRangeDate) ([]*entity.Traffic, error) {
	var traffic []*entity.Traffic
	err := repo.db.WithContext(ctx).
		Model(&entity.Traffic{}).
		Select("date, SUM(\"in\") AS \"in\", SUM(out) AS out, SUM(bandwidth) as bandwidth, SUM(bytes_in) AS bytes_in, SUM(bytes_out) AS bytes_out").
		Joins("JOIN fat_interfaces ON fat_interfaces.interface_id = traffics.interface_id").
		Joins("JOIN fats ON fats.id = fat_interfaces.fat_id").
		Joins("JOIN locations as l ON l.id = fats.location_id").
		Where("l.state = ? AND traffics.date BETWEEN ? AND ?", state, date.InitDate, date.EndDate).
		Group("date").
		Order("date").
		Find(&traffic).
		Error

	return traffic, err
}

func (repo trafficRepository) GetTrafficByCounty(ctx context.Context, state, county string, date *model.TranficRangeDate) ([]*entity.Traffic, error) {
	var traffic []*entity.Traffic
	err := repo.db.WithContext(ctx).
		Model(&entity.Traffic{}).
		Select("date, SUM(\"in\") AS \"in\", SUM(out) AS out, SUM(bandwidth) as bandwidth, SUM(bytes_in) AS bytes_in, SUM(bytes_out) AS bytes_out").
		Joins("JOIN fat_interfaces ON fat_interfaces.interface_id = traffics.interface_id").
		Joins("JOIN fats ON fats.id = fat_interfaces.fat_id").
		Joins("JOIN locations as l ON l.id = fats.location_id").
		Where("l.state = ? AND l.county = ? AND traffics.date BETWEEN ? AND ?", state, county, date.InitDate, date.EndDate).
		Group("date").
		Order("date").
		Find(&traffic).
		Error

	return traffic, err
}

func (repo trafficRepository) GetTrafficByMunicipality(ctx context.Context, state, county, municipality string, date *model.TranficRangeDate) ([]*entity.Traffic, error) {
	var traffic []*entity.Traffic
	err := repo.db.WithContext(ctx).
		Model(&entity.Traffic{}).
		Select("date, SUM(\"in\") AS \"in\", SUM(out) AS out, SUM(bandwidth) as bandwidth, SUM(bytes_in) AS bytes_in, SUM(bytes_out) AS bytes_out").
		Joins("JOIN fat_interfaces ON fat_interfaces.interface_id = traffics.interface_id").
		Joins("JOIN fats ON fats.id = fat_interfaces.fat_id").
		Joins("JOIN locations as l ON l.id = fats.location_id").
		Where("l.state = ? AND l.county = ? AND l.municipality = ? AND traffics.date BETWEEN ? AND ?", state, county, municipality, date.InitDate, date.EndDate).
		Group("date").
		Order("date").
		Find(&traffic).
		Error

	return traffic, err
}

func (repo trafficRepository) GetTrafficByODN(ctx context.Context, odn string, date *model.TranficRangeDate) ([]*entity.Traffic, error) {
	var traffic []*entity.Traffic
	err := repo.db.WithContext(ctx).
		Model(&entity.Traffic{}).
		Select("date, SUM(\"in\") AS \"in\", SUM(out) AS out, SUM(bandwidth) as bandwidth, SUM(bytes_in) AS bytes_in, SUM(bytes_out) AS bytes_out").
		Joins("JOIN fat_interfaces ON fat_interfaces.interface_id = traffics.interface_id").
		Joins("JOIN fats ON fats.id = fat_interfaces.fat_id").
		Where("fats.odn = ? AND traffics.date BETWEEN ? AND ?", odn, date.InitDate, date.EndDate).
		Group("date").
		Order("date").
		Find(&traffic).
		Error

	return traffic, err
}

func (repo trafficRepository) GetTotalTrafficByState(ctx context.Context, ids []uint, month string) (*model.TrafficState, error) {
	var trafficByState *model.TrafficState
	err := repo.db.WithContext(ctx).
		Model(&entity.Trend{}).
		Select("SUM(bandwidth) AS bandwidth, SUM(\"in\") AS \"in\", SUM(out) AS out, SUM(bytes_in) AS bytes_in, SUM(bytes_out) AS bytes_out").
		Where("device_id IN (?) AND Extract(MONTH FROM date) = ?", ids, month).
		Scan(&trafficByState).
		Error
	return trafficByState, err
}

func (repo trafficRepository) GetTotalTrafficByODN(ctx context.Context, ids []uint, month string) (*model.TrafficODN, error) {
	var trafficByState *model.TrafficODN
	err := repo.db.WithContext(ctx).
		Model(&entity.Trend{}).
		Select("SUM(bandwidth) AS bandwidth, SUM(\"in\") AS \"in\", SUM(out) AS out, SUM(bytes_in) AS bytes_in, SUM(bytes_out) AS bytes_out").
		Where("device_id IN (?) AND Extract(MONTH FROM date) = ?", ids, month).
		Scan(&trafficByState).
		Error
	return trafficByState, err
}
