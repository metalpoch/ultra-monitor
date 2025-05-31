package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type FatRepository interface {
	AddFat(ctx context.Context, tao *entity.Fat) error
	DeleteOne(ctx context.Context, id int32) error
	AllTao(ctx context.Context, page, limit uint16) ([]entity.Fat, error)
	GetByID(ctx context.Context, id int32) (entity.Fat, error)
	GetByFat(ctx context.Context, tao string) (entity.Fat, error)
	GetByOdn(ctx context.Context, state, odn string) ([]entity.Fat, error)
	GetOdnByStates(ctx context.Context, state string) ([]string, error)
	GetOdnByCounty(ctx context.Context, state, county string) ([]string, error)
	GetOdnMunicipality(ctx context.Context, state, county, municipality string) ([]string, error)
	GetOdnByOlt(ctx context.Context, oltIP string) ([]string, error)
	GetOdnByOltPort(ctx context.Context, oltIP string, shell, card, port uint8) ([]string, error)
	GetAllOdn(ctx context.Context) ([]string, error)
}

type fatRepository struct {
	db *sqlx.DB
}

func NewFatRepository(db *sqlx.DB) *fatRepository {
	return &fatRepository{db}
}

func (repo *fatRepository) AddFat(ctx context.Context, tao *entity.Fat) error {
	query := `
        INSERT INTO fats (
            region, fat, state, municipality, county, odn, olt_ip,
            pon_shell, pon_port, pon_card, latitude, longitude
        ) VALUES (
            :region, :fat, :state, :municipality, :county, :sector, :odn, :olt_ip,
            :pon_shell, :pon_port, :pon_card, :latitude, :longitude
        )`
	_, err := repo.db.NamedExecContext(ctx, query, tao)
	return err
}

func (repo *fatRepository) DeleteOne(ctx context.Context, id int32) error {
	query := `DELETE FROM fats WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *fatRepository) AllTao(ctx context.Context, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	query := `SELECT * FROM fats ORDER BY region, state, municipality, county LIMIT $1 OFFSET $2`
	err := repo.db.SelectContext(ctx, &res, query, limit, offset)
	return res, err
}

func (repo *fatRepository) GetByID(ctx context.Context, id int32) (entity.Fat, error) {
	var res entity.Fat
	query := `SELECT * FROM fats WHERE id = $1`
	err := repo.db.GetContext(ctx, &res, query, id)
	return res, err
}

func (repo *fatRepository) GetByFat(ctx context.Context, tao string) (entity.Fat, error) {
	var res entity.Fat
	query := `SELECT * FROM fats WHERE fat = $1`
	err := repo.db.SelectContext(ctx, &res, query, tao)
	return res, err
}

func (repo *fatRepository) GetByOdn(ctx context.Context, state, odn string) ([]entity.Fat, error) {
	var res []entity.Fat
	query := `SELECT * FROM fats WHERE state = $1 AND odn = $2 ORDER BY fat`
	err := repo.db.SelectContext(ctx, &res, query, state, odn)
	return res, err
}

func (repo *fatRepository) GetOdnByStates(ctx context.Context, state string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM fats WHERE state = $1 ORDER BY odn`
	err := repo.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (repo *fatRepository) GetOdnByCounty(ctx context.Context, state, county string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM fats WHERE state = $1 AND county = $2 ORDER BY odn`
	err := repo.db.SelectContext(ctx, &res, query, state, county)
	return res, err
}

func (repo *fatRepository) GetOdnMunicipality(ctx context.Context, state, county, municipality string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM fats WHERE state = $1 AND county = $2 AND municipality = $3 ORDER BY odn`
	err := repo.db.SelectContext(ctx, &res, query, state, county, municipality)
	return res, err
}

func (repo *fatRepository) GetOdnByOlt(ctx context.Context, oltIP string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM fats WHERE olt_ip = $1 ORDER BY odn`
	err := repo.db.SelectContext(ctx, &res, query, oltIP)
	return res, err
}

func (repo *fatRepository) GetOdnByOltPort(ctx context.Context, oltIP string, shell, card, port uint8) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM fats WHERE olt_ip = $1 AND shell = $2 AND card = $3 AND port = $4 ORDER BY odn`
	err := repo.db.SelectContext(ctx, &res, query, oltIP, shell, card, port)
	return res, err
}

func (repo *fatRepository) GetAllOdn(ctx context.Context) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM fats`
	err := repo.db.SelectContext(ctx, &res, query)
	return res, err
}
