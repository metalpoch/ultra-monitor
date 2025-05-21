package repository

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"gorm.io/gorm"
)

type SummaryRepository interface {
	UserStatus(ctx context.Context, initDate, endDate time.Time) ([]entity.UserStatusCounts, error)
}

type summaryRepository struct {
	db *gorm.DB
}

func NewSummaryRepository(db *gorm.DB) *summaryRepository {
	return &summaryRepository{db}
}

func (repo summaryRepository) UserStatus(ctx context.Context, initDate, endDate time.Time) ([]entity.UserStatusCounts, error) {
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
		initDate,
		endDate,
	).Scan(&results).Error
	return results, err
}

/*
 * SELECT DATE_TRUNC('hour', date) AS hour, COUNT(*) AS count
 *FROM measurement_onts
 *WHERE control_run_status = 1
 *GROUP BY hour
 *ORDER BY hour;*/

/*        err := repo.db.WithContext(ctx).Table("measurement_onts").
 *                Select("DATE_TRUNC('hour', date) AS hour, COUNT(*) AS count").
 *                Where("control_run_status = ?", 1).
 *                Group("hour").
 *                Order("hour").
 *                Find(&results).
 *                Error
 *
 *        return results, err*/
