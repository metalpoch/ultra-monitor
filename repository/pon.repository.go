package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type PonRepository interface {
	PonByPort(ctx context.Context, sysname, port string) (entity.Pon, error)
	PonsByOLT(ctx context.Context, sysname string) ([]entity.Pon, error)
}

type ponRepository struct {
	db *sqlx.DB
}

func NewPonRepository(db *sqlx.DB) *ponRepository {
	return &ponRepository{db}
}

func (repo *ponRepository) PonByPort(ctx context.Context, sysname, port string) (entity.Pon, error) {
	var res entity.Pon
	query := `
	SELECT id, olt_id, if_index, if_name, if_descr, if_alias, created_at, updated_at
	FROM pon
	JOIN olt ON olt.id = pon.olt_id
	WHERE olt.sys_name = ? AND pon.if_name = ?`
	err := repo.db.SelectContext(ctx, &res, query, sysname, port)
	return res, err
}

func (repo *ponRepository) PonsByOLT(ctx context.Context, sysname string) ([]entity.Pon, error) {
	var res []entity.Pon
	query := `
	SELECT id, olt_id, if_index, if_name, if_descr, if_alias, created_at, updated_at
	FROM pon
	JOIN olt ON olt.id = pon.olt_id
	WHERE olt.sys_name = ?`
	err := repo.db.SelectContext(ctx, &res, query, sysname)
	return res, err
}
