package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/entity"
)

type TaoRepository interface {
	AddTao(ctx context.Context, tao *entity.Tao) error
	DeleteOne(ctx context.Context, id uint64) error
	AllTao(ctx context.Context, page, limit uint16) ([]entity.Tao, error)
	GetByID(ctx context.Context, id uint64) (entity.Tao, error)
	GetByTao(ctx context.Context, tao string) (entity.Tao, error)
	GetByOdn(ctx context.Context, state, odn string) ([]entity.Tao, error)
	GetOdnStates(ctx context.Context, state string) ([]string, error)
	GetOdnCounty(ctx context.Context, state, county string) ([]string, error)
	GetOdnMunicipality(ctx context.Context, state, county, municipality string) ([]string, error)
	GetOdnByOlt(ctx context.Context, oltIP string) ([]string, error)
	GetOdnOltPort(ctx context.Context, oltIP string, shell, card, port uint8) ([]string, error)
	GetAllOdn(ctx context.Context) ([]string, error)
}

type taoRepository struct {
	db *sqlx.DB
}

func NewTaoRepository(db *sqlx.DB) *taoRepository {
	return &taoRepository{db}
}

func (repo *taoRepository) AddTao(ctx context.Context, tao *entity.Tao) error {
	query := `
        INSERT INTO tao (
            region, state, municipality, county, sector, odn, olt_ip,
            pon_shell, pon_port, pon_card, latitude, longitude
        ) VALUES (
            :region, :state, :municipality, :county, :sector, :odn, :olt_ip,
            :pon_shell, :pon_port, :pon_card, :latitude, :longitude
        )
    `
	_, err := repo.db.NamedExecContext(ctx, query, tao)
	return err
}

func (repo *taoRepository) DeleteOne(ctx context.Context, id uint64) error {
	query := `DELETE FROM tao WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *taoRepository) AllTao(ctx context.Context, page, limit uint16) ([]entity.Tao, error) {
	var res []entity.Tao
	offset := (page - 1) * limit
	query := `SELECT * FROM tao LIMIT $1 OFFSET $2`
	err := repo.db.SelectContext(ctx, &res, query, limit, offset)
	return res, err
}

func (repo *taoRepository) GetByID(ctx context.Context, id uint64) (entity.Tao, error) {
	var res entity.Tao
	query := `SELECT * FROM tao WHERE id = $1`
	err := repo.db.GetContext(ctx, &res, query, id)
	return res, err
}

func (repo *taoRepository) GetByTao(ctx context.Context, tao string) (entity.Tao, error) {
	var res entity.Tao
	query := `SELECT * FROM tao WHERE tao = $1`
	err := repo.db.SelectContext(ctx, &res, query, tao)
	return res, err
}

func (repo *taoRepository) GetByOdn(ctx context.Context, state, odn string) ([]entity.Tao, error) {
	var res []entity.Tao
	query := `SELECT * FROM tao WHERE state = $1 AND odn = $2`
	err := repo.db.SelectContext(ctx, &res, query, state, odn)
	return res, err
}

func (repo *taoRepository) GetOdnStates(ctx context.Context, state string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM tao WHERE state = $1`
	err := repo.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (repo *taoRepository) GetOdnCounty(ctx context.Context, state, county string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM tao WHERE state = $1 AND county = $2`
	err := repo.db.SelectContext(ctx, &res, query, state, county)
	return res, err
}

func (repo *taoRepository) GetOdnMunicipality(ctx context.Context, state, county, municipality string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM tao WHERE state = $1 AND county = $2 AND municipality = $3`
	err := repo.db.SelectContext(ctx, &res, query, state, county, municipality)
	return res, err
}

func (repo *taoRepository) GetOdnByOlt(ctx context.Context, oltIP string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM tao WHERE olt_ip = $1`
	err := repo.db.SelectContext(ctx, &res, query, oltIP)
	return res, err
}

func (repo *taoRepository) GetOdnOltPort(ctx context.Context, oltIP string, shell, card, port uint8) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM tao WHERE olt_ip = $1 AND shell = $2 AND card = $3 AND port = $4`
	err := repo.db.SelectContext(ctx, &res, query, oltIP, shell, card, port)
	return res, err
}

func (repo *taoRepository) GetAllOdn(ctx context.Context) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT odn FROM tao`
	err := repo.db.SelectContext(ctx, &res, query)
	return res, err
}
