package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/entity"
)

type LocationRepository interface {
	States(ctx context.Context) ([]string, error)
	Counties(ctx context.Context, state string) ([]string, error)
	Municipalities(ctx context.Context, state, county string) ([]string, error)
}

type locationRepository struct {
	db *sqlx.DB
}

func NewLocationRepository(db *sqlx.DB) *locationRepository {
	return &locationRepository{db}
}

func (repo locationRepository) Add(ctx context.Context, location entity.Location) error {
	query := `INSERT INTO location (state, county, municipality) VALUES (:state, :county, :municipality)`
	_, err := repo.db.NamedExecContext(ctx, query, location)
	return err
}

func (repo locationRepository) States(ctx context.Context) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT state FROM location ORDER BY state`
	err := repo.db.SelectContext(ctx, &res, query)
	return res, err
}

func (repo locationRepository) Counties(ctx context.Context, state string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT county FROM location WHERE state = ? ORDER BY state`
	err := repo.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (repo locationRepository) Municipalities(ctx context.Context, state, county string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT municipality FROM location WHERE state = ? AND county = ? ORDER BY state`
	err := repo.db.SelectContext(ctx, &res, query, state)
	return res, err
}
