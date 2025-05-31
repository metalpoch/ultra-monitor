package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/constants"
)

type PonRepository interface {
	PonsByOLT(ctx context.Context, sysname string) ([]entity.Pon, error)
	PonByPort(ctx context.Context, sysname, port string) (entity.Pon, error)
}

type ponRepository struct {
	db *sqlx.DB
}

func NewPonRepository(db *sqlx.DB) *ponRepository {
	return &ponRepository{db}
}

func (repo *ponRepository) PonsByOLT(ctx context.Context, sysname string) ([]entity.Pon, error) {
	var res []entity.Pon
	err := repo.db.SelectContext(ctx, &res, constants.SQL_PONS_BY_OLT, sysname)
	return res, err
}
func (repo *ponRepository) PonByPort(ctx context.Context, sysname, port string) (entity.Pon, error) {
	var res entity.Pon
	err := repo.db.SelectContext(ctx, &res, constants.SQL_PON_BY_PORT, sysname, port)
	return res, err
}
