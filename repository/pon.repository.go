package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
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
	query := `SELECT pons.* FROM pons JOIN olts ON olts.ip = pons.olt_ip WHERE olts.sys_name = $1`
	err := repo.db.SelectContext(ctx, &res, query, sysname)
	return res, err
}

func (repo *ponRepository) PonByPort(ctx context.Context, sysname, port string) (entity.Pon, error) {
	var res entity.Pon
	query := `SELECT pons.* FROM pons JOIN olts ON olts.ip = pons.olt_ip WHERE olts.sys_name = $1 AND pons.if_name = $2`
	err := repo.db.GetContext(ctx, &res, query, sysname, port)
	return res, err
}
