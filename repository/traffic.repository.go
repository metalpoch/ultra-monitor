package repository

import (
	"github.com/jmoiron/sqlx"
)

type TrafficRepository interface {
}

type trafficRepository struct {
	db *sqlx.DB
}

func NewTrafficRepository(db *sqlx.DB) *trafficRepository {
	return &trafficRepository{db}
}

func (r *trafficRepository) keloke() {}
