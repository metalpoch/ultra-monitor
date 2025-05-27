package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/entity"
)

type OdnRepository interface {
	GetFat(ctx context.Context, id uint) (entity.Fat, error)
	GetODN(ctx context.Context, odn string) ([]entity.FatInterface, error)
	GetODNStates(ctx context.Context, state string) ([]string, error)
	GetODNStatesContries(ctx context.Context, state, country string) ([]string, error)
	GetODNStatesContriesMunicipality(ctx context.Context, state, country, municipality string) ([]string, error)
	GetODNDevice(ctx context.Context, id uint) ([]string, error)
	GetODNDevicePort(ctx context.Context, id uint, pattern string) ([]string, error)
	GetAllODN(ctx context.Context) ([]string, error)
}

type odnRepository struct {
	db *sqlx.DB
}

func NewOdnRepository(db *sqlx.DB) *odnRepository {
	return &odnRepository{db}
}

func (repo odnRepository) GetFat(ctx context.Context, id uint) (entity.Fat, error) {
	var res entity.Fat
	query := `SELECT * FROM fats WHERE id = ?`
	err := repo.db.GetContext(ctx, &res, query, id)
	return res, err
}

func (repo odnRepository) GetODN(ctx context.Context, odn string) ([]entity.FatInterface, error) {
	var res []entity.FatInterface
	query := `
        SELECT fi.* FROM fats_pon fi
        JOIN fats f ON f.id = fi.fat_id
        WHERE f.odn = ?
    `
	err := repo.db.SelectContext(ctx, &res, query, odn)
	return res, err
}

func (repo odnRepository) GetODNStates(ctx context.Context, state string) ([]string, error) {
	var res []string
	query := `
        SELECT DISTINCT f.odn
        FROM fats f
        JOIN locations l ON f.location_id = l.id
        WHERE l.state = ? `
	err := repo.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (repo odnRepository) GetODNStatesContries(ctx context.Context, state, country, municipality string) ([]string, error) {
	var res []string
	query := `
        SELECT DISTINCT f.odn
        FROM fats f
        JOIN locations l ON f.location_id = l.id
        WHERE l.state = ? AND l.county = ?
    `
	err := repo.db.SelectContext(ctx, &res, query, state, country, municipality)
	return res, err
}

func (repo odnRepository) GetODNStatesContriesMunicipality(ctx context.Context, state, country, municipality string) ([]string, error) {
	var res []string
	query := `
        SELECT DISTINCT f.odn
        FROM fats f
        JOIN locations l ON f.location_id = l.id
        WHERE l.state = ? AND l.county = ? AND l.municipality = ?
    `
	err := repo.db.SelectContext(ctx, &res, query, state, country, municipality)
	return res, err
}

func (repo odnRepository) GetODNDevice(ctx context.Context, id uint) ([]string, error) {
	var res []string
	query := `
        SELECT DISTINCT f.odn
        FROM pons p
        INNER JOIN fats_pon fp ON p.id = fp.interface_id
        INNER JOIN fats f ON fp.fat_id = f.id
        WHERE p.device_id = ?
    `
	err := repo.db.SelectContext(ctx, &res, query, id)
	return res, err
}

func (repo odnRepository) GetODNDevicePort(ctx context.Context, id uint, pattern string) ([]string, error) {
	var res []string
	query := `
        SELECT DISTINCT f.odn
        FROM pons p
        INNER JOIN fats_pon fp ON p.id = fp.interface_id
        INNER JOIN fats f ON fp.fat_id = f.id
        WHERE p.device_id = ? AND p.if_name LIKE ?
    `
	err := repo.db.SelectContext(ctx, &res, query, id, pattern)
	return res, err
}

func (repo odnRepository) GetAllODN(ctx context.Context) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM fats`
	err := repo.db.SelectContext(ctx, res, query)
	return res, err
}

func (repo odnRepository) GetDevicesByOND(ctx context.Context, odn string) ([]uint, error) {
	var res []uint
	query := `
        SELECT DISTINCT d.id
        FROM devices d
        INNER JOIN pons p ON d.id = p.device_id
        INNER JOIN fats_pon fp ON p.id = fp.interface_id
        INNER JOIN fats f ON fp.fat_id = f.id
        WHERE f.odn = ?
    `
	err := repo.db.SelectContext(ctx, &res, query, odn)
	return res, err
}
