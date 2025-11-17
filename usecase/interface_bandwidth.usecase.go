package usecase

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/mongodb"
	"github.com/metalpoch/ultra-monitor/repository"
)

type InterfaceBandwidthUsecase struct {
	repo repository.InterfaceBandwidthRepository
}

func NewInterfaceBandwidthUsecase(db *sqlx.DB, mongoDB *mongodb.MongoDB) *InterfaceBandwidthUsecase {
	return &InterfaceBandwidthUsecase{
		repo: repository.NewInterfaceBandwidthRepository(db, mongoDB.Database),
	}
}

func (u *InterfaceBandwidthUsecase) UpdateInterfaceBandwidth(ctx context.Context) error {
	// Calculate date range: yesterday 00:00 to today 23:59
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	startDate := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	// Get interface bandwidth data from MongoDB
	bandwidthData, err := u.repo.GetInterfaceBandwidthFromMongoDB(ctx, startDate, endDate)
	if err != nil {
		return err
	}

	if len(bandwidthData) == 0 {
		return nil // No data to process
	}

	// Clean previously data
	err = u.repo.CleanInterfacesBandwidth(ctx)
	if err != nil {
		return err
	}

	// Store the data in PostgreSQL
	err = u.repo.InsertInterfaceBandwidth(ctx, bandwidthData)
	if err != nil {
		return err
	}

	return nil
}
