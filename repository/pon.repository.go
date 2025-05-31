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
	query := `
	SELECT pons.*
	FROM pons
	JOIN olts ON olts.id = pons.olt_id
	WHERE olts.sys_name = ?`
	err := repo.db.SelectContext(ctx, &res, query, sysname)
	return res, err
}
func (repo *ponRepository) PonByPort(ctx context.Context, sysname, port string) (entity.Pon, error) {
	var res entity.Pon
	query := `
	SELECT pons.*
	FROM pons
	JOIN olts ON olts.id = pons.olt_id
	WHERE olts.sys_name = ? AND pons.if_name = ?`
	err := repo.db.SelectContext(ctx, &res, query, sysname, port)
	return res, err
}
