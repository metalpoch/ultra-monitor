package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
	commonModel "github.com/metalpoch/olt-blueprint/common/model"
	model "github.com/metalpoch/olt-blueprint/core/model"
	"gorm.io/gorm"
)

type SummaryRepository interface {
	UserStatus(ctx context.Context, date *commonModel.TranficRangeDate) ([]entity.UserStatusCounts, error)
	UserStatusByState(ctx context.Context, query *model.UserStatusByState) ([]entity.UserStatusCounts, error)
	Traffic(ctx context.Context, date *commonModel.TranficRangeDate) ([]entity.TrafficResponse, error)
	TrafficByState(ctx context.Context, state string, date *commonModel.TranficRangeDate) ([]entity.TrafficResponse, error)
}

type summaryRepository struct {
	db *gorm.DB
}

func NewSummaryRepository(db *gorm.DB) *summaryRepository {
	return &summaryRepository{db}
}

func (repo summaryRepository) UserStatus(ctx context.Context, date *commonModel.TranficRangeDate) ([]entity.UserStatusCounts, error) {
	var results []entity.UserStatusCounts
	err := repo.db.WithContext(ctx).Raw(`
    		SELECT 
			DATE_TRUNC('hour', date) AS hour, 
			COUNT(*) FILTER (WHERE control_run_status = 1) AS active_count,
			COUNT(*) FILTER (WHERE control_run_status = 2) AS inactive_count,
			COUNT(*) FILTER (WHERE control_run_status NOT IN (1,2)) AS unknown_count,
			COUNT(*) AS total_count
		FROM measurement_onts
		WHERE date BETWEEN ? AND ?
		GROUP BY hour
		ORDER BY hour`,
		date.InitDate,
		date.EndDate,
	).Scan(&results).Error
	return results, err
}

func (repo summaryRepository) UserStatusByState(ctx context.Context, query *model.UserStatusByState) ([]entity.UserStatusCounts, error) {
	var results []entity.UserStatusCounts
	err := repo.db.WithContext(ctx).Raw(`
    		SELECT
        		DATE_TRUNC('hour', measurement_onts.date) AS hour,
			COUNT(CASE WHEN control_run_status = 1 THEN 1 END) AS active_count,
			COUNT(CASE WHEN control_run_status = 2 THEN 1 END) AS inactive_count,
			COUNT(CASE WHEN control_run_status NOT IN (1,2) THEN 1 END) AS unknown_count,
			COUNT(*) AS total_count
		FROM measurement_onts
		INNER JOIN fat_interfaces ON fat_interfaces.interface_id = measurement_onts.interface_id
		INNER JOIN fats ON fats.id = fat_interfaces.fat_id
		INNER JOIN locations AS l ON l.id = fats.location_id
		WHERE l.state IN (?) AND measurement_onts.date >= ? AND measurement_onts.date < ?
		GROUP BY hour
		ORDER BY hour`,
		query.States,
		query.InitDate,
		query.EndDate,
	).Scan(&results).Error
	return results, err
}

func (repo summaryRepository) Traffic(ctx context.Context, date *commonModel.TranficRangeDate) ([]entity.TrafficResponse, error) {
	var trafficTrend []entity.TrafficResponse
	err := repo.db.WithContext(ctx).
		Model(&entity.Traffic{}).
		Select("DATE_TRUNC('hour', date) AS date",
			"SUM(\"in\") / 1000000 AS mbps_in",
			"SUM(out) / 1000000 AS mbps_out",
			"SUM(bandwidth) / 1000000 AS bandwidth",
			"SUM(bytes_in) / 1000000 AS mbytes_in",
			"SUM(bytes_out) /1000000 AS mbytes_out").
		Where("date BETWEEN ? AND ?", date.InitDate, date.EndDate).
		Group("DATE_TRUNC('hour', date)").
		Order("date").
		Find(&trafficTrend).
		Error

	return trafficTrend, err
}

func (repo summaryRepository) TrafficByState(ctx context.Context, state string, date *commonModel.TranficRangeDate) ([]entity.TrafficResponse, error) {
	var trafficTrend []entity.TrafficResponse
	err := repo.db.WithContext(ctx).
		Model(&entity.Traffic{}).
		Select("DATE_TRUNC('hour', traffics.date) AS date",
			"SUM(\"in\") AS \"in\"",
			"SUM(out) AS out",
			"SUM(bandwidth) AS bandwidth",
			"SUM(bytes_in) AS bytes_in",
			"SUM(bytes_out) AS bytes_out").
		Joins("JOIN fat_interfaces ON fat_interfaces.interface_id = traffics.interface_id").
		Joins("JOIN fats ON fats.id = fat_interfaces.fat_id").
		Joins("JOIN locations AS l ON l.id = fats.location_id").
		Where("l.state = ?", state).
		Where("traffics.date BETWEEN ? AND ?", date.InitDate, date.EndDate).
		Group("DATE_TRUNC('hour', traffics.date)").
		Order("date").
		Find(&trafficTrend).
		Error

	return trafficTrend, err
}
