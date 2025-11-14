package repository

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type InterfaceBandwidthRepository interface {
	// MongoDB operations
	GetInterfaceBandwidthFromMongoDB(ctx context.Context, startDate, endDate time.Time) ([]entity.InterfaceBandwidth, error)

	// PostgreSQL operations
	UpsertInterfaceBandwidth(ctx context.Context, bandwidthData []entity.InterfaceBandwidth) error
}

type interfaceBandwidthRepository struct {
	db *sqlx.DB
}

func NewInterfaceBandwidthRepository(db *sqlx.DB) *interfaceBandwidthRepository {
	return &interfaceBandwidthRepository{db}
}

func (r *interfaceBandwidthRepository) GetInterfaceBandwidthFromMongoDB(ctx context.Context, startDate, endDate time.Time) ([]entity.InterfaceBandwidth, error) {
	// This method will be implemented in a separate MongoDB-specific repository
	// For now, return empty slice - the actual implementation will be in the MongoDB repository
	return []entity.InterfaceBandwidth{}, nil
}

func (r *interfaceBandwidthRepository) UpsertInterfaceBandwidth(ctx context.Context, bandwidthData []entity.InterfaceBandwidth) error {
	query := `
		INSERT INTO interface_bandwidth (interface, olt, bandwidth, created_at)
		VALUES (:interface, :olt, :bandwidth, :created_at)
		ON CONFLICT (interface, created_at) DO UPDATE SET
			olt = EXCLUDED.olt,
			bandwidth = EXCLUDED.bandwidth
	`

	_, err := r.db.NamedExecContext(ctx, query, bandwidthData)
	return err
}

