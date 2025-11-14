package usecase

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/mongodb"
	"github.com/metalpoch/ultra-monitor/repository"
)

type InterfaceBandwidthUsecase struct {
	postgresRepo repository.InterfaceBandwidthRepository
	mongoRepo    *repository.InterfaceBandwidthMongoRepository
}

func NewInterfaceBandwidthUsecase(db *sqlx.DB, mongoDB *mongodb.MongoDB) *InterfaceBandwidthUsecase {
	return &InterfaceBandwidthUsecase{
		postgresRepo: repository.NewInterfaceBandwidthRepository(db),
		mongoRepo:    repository.NewInterfaceBandwidthMongoRepository(mongoDB.Database),
	}
}

func (u *InterfaceBandwidthUsecase) UpdateInterfaceBandwidth(ctx context.Context) error {
	// Calculate date range: yesterday 00:00 to today 23:59
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)

	startDate := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endDate := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	// Get interface bandwidth data from MongoDB
	bandwidthData, err := u.mongoRepo.GetInterfaceBandwidthFromMongoDB(ctx, startDate, endDate)
	if err != nil {
		return err
	}

	if len(bandwidthData) == 0 {
		return nil // No data to process
	}

	// Store the data in PostgreSQL
	err = u.postgresRepo.UpsertInterfaceBandwidth(ctx, bandwidthData)
	if err != nil {
		return err
	}

	return nil
}

