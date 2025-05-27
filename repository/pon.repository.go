package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/entity"
)

type PonRepository interface {
	Upsert(ctx context.Context, element entity.Pon) (uint64, error)
	GetPon(ctx context.Context, id uint64) (entity.Pon, error)
	PonByPort(ctx context.Context, sysname, pattern string) ([]entity.Pon, error)
	PonsByOLT(ctx context.Context, sysname string) ([]entity.Pon, error)
}

type ponRepository struct {
	db *sqlx.DB
}

func NewPonRepository(db *sqlx.DB) *ponRepository {
	return &ponRepository{db}
}

func (repo ponRepository) Upsert(ctx context.Context, element entity.Pon) (uint64, error) {
	var id uint64
	query := `
        INSERT INTO pon (olt_id, if_index, if_name, if_descr, if_alias, created_at, updated_at)
        VALUES (:olt_id, :if_index, :if_name, :if_descr, :if_alias, :created_at, :updated_at)
        ON CONFLICT (olt_id, if_index) DO UPDATE SET
            if_name = EXCLUDED.if_name,
            if_descr = EXCLUDED.if_descr,
            if_alias = EXCLUDED.if_alias,
            updated_at = EXCLUDED.updated_at
        RETURNING id
    `
	err := repo.db.QueryRowxContext(ctx, query, element).Scan(&id)
	return id, err
}

func (repo ponRepository) GetPon(ctx context.Context, id uint64) (entity.Pon, error) {
	var res entity.Pon
	query := `SELECT * FROM pon WHERE id = ?`
	err := repo.db.SelectContext(ctx, &res, query, id)
	return res, err
}

func (repo ponRepository) PonByPort(ctx context.Context, sysname, port string) (entity.Pon, error) {
	var res entity.Pon
	query := `
	SELECT id, olt_id, if_index, if_name, if_descr, if_alias, created_at, updated_at
	FROM pon
	JOIN olt ON olt.id = pon.olt_id
	WHERE olt.sys_name = ? AND pon.if_name = ?`
	err := repo.db.SelectContext(ctx, &res, query, sysname, port)
	return res, err
}

func (repo ponRepository) PonsByOLT(ctx context.Context, sysname string) ([]entity.Pon, error) {
	var res []entity.Pon
	query := `
	SELECT id, olt_id, if_index, if_name, if_descr, if_alias, created_at, updated_at
	FROM pon
	JOIN olt ON olt.id = pon.olt_id
	WHERE olt.sys_name = ?`
	err := repo.db.SelectContext(ctx, &res, query, sysname)
	return res, err
}
