package repository

import "github.com/jmoiron/sqlx"

type OntRepository interface {
	GetStatusTrend(id int64) (any, error)
}

type ontRepository struct {
	repo *sqlx.DB
}

func NewOntRepository(db *sqlx.DB) *ontRepository {
	return &ontRepository{db}
}

func (r *ontRepository) GetStatusTrend(id int64) (any, error) {
	return nil, nil
}
