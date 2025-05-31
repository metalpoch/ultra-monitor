package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/constants"
)

type FatRepository interface {
	AddFat(ctx context.Context, tao *entity.Fat) error
	DeleteOne(ctx context.Context, id int32) error
	AllFat(ctx context.Context, page, limit uint16) ([]entity.Fat, error)
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
	_, err := repo.db.NamedExecContext(ctx, constants.SQL_INSERT_FAT, tao)
	return err
}

func (repo *fatRepository) DeleteOne(ctx context.Context, id int32) error {
	_, err := repo.db.ExecContext(ctx, constants.SQL_DELETE_FAT_BY_ID, id)
	return err
}

func (repo *fatRepository) AllFat(ctx context.Context, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	err := repo.db.SelectContext(ctx, &res, constants.SQL_SELECT_ALL_FATS, limit, offset)
	return res, err
}

func (repo *fatRepository) GetByID(ctx context.Context, id int32) (entity.Fat, error) {
	var res entity.Fat
	err := repo.db.GetContext(ctx, &res, constants.SQL_SELECT_FAT_BY_ID, id)
	return res, err
}

func (repo *fatRepository) GetByFat(ctx context.Context, tao string) (entity.Fat, error) {
	var res entity.Fat
	err := repo.db.SelectContext(ctx, &res, constants.SQL_SELECT_FAT_BY_FAT, tao)
	return res, err
}

func (repo *fatRepository) GetByOdn(ctx context.Context, state, odn string) ([]entity.Fat, error) {
	var res []entity.Fat
	err := repo.db.SelectContext(ctx, &res, constants.SQL_SELECT_FATS_BY_ODN, state, odn)
	return res, err
}

func (repo *fatRepository) GetOdnByStates(ctx context.Context, state string) ([]string, error) {
	var res []string
	err := repo.db.SelectContext(ctx, &res, constants.SQL_SELECT_DISTINCT_ODN_BY_STATE, state)
	return res, err
}

func (repo *fatRepository) GetOdnByCounty(ctx context.Context, state, county string) ([]string, error) {
	var res []string
	err := repo.db.SelectContext(ctx, &res, constants.SQL_SELECT_DISTINCT_ODN_BY_COUNTY, state, county)
	return res, err
}

func (repo *fatRepository) GetOdnMunicipality(ctx context.Context, state, county, municipality string) ([]string, error) {
	var res []string
	err := repo.db.SelectContext(ctx, &res, constants.SQL_SELECT_DISTINCT_ODN_BY_MUNICIPALITY, state, county, municipality)
	return res, err
}

func (repo *fatRepository) GetOdnByOlt(ctx context.Context, oltIP string) ([]string, error) {
	var res []string
	err := repo.db.SelectContext(ctx, &res, constants.SQL_SELECT_DISTINCT_ODN_BY_OLT, oltIP)
	return res, err
}

func (repo *fatRepository) GetOdnByOltPort(ctx context.Context, oltIP string, shell, card, port uint8) ([]string, error) {
	var res []string
	err := repo.db.SelectContext(ctx, &res, constants.SQL_SELECT_DISTINCT_ODN_BY_OLT_PORT, oltIP, shell, card, port)
	return res, err
}

func (repo *fatRepository) GetAllOdn(ctx context.Context) ([]string, error) {
	var res []string
	err := repo.db.SelectContext(ctx, &res, constants.SQL_SELECT_DISTINCT_ALL_ODN)
	return res, err
}
