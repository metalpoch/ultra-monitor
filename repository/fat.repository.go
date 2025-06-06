package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/constants"
)

type FatRepository interface {
	AllFat(ctx context.Context, page, limit uint16) ([]entity.Fat, error)
	AddFat(ctx context.Context, tao *entity.Fat) error
	DeleteOne(ctx context.Context, id int32) error
	GetByID(ctx context.Context, id int32) (entity.Fat, error)
	GetStates(ctx context.Context) ([]string, error)
	GetMunicipality(ctx context.Context, state string) ([]string, error)
	GetCounty(ctx context.Context, state, municipality string) ([]string, error)
	GetFatsByStates(ctx context.Context, state string) ([]entity.Fat, error)
	GetFatsByMunicipality(ctx context.Context, state, municipality string) ([]entity.Fat, error)
	GetFatsByCounty(ctx context.Context, state, municipality, county string) ([]entity.Fat, error)
}

type fatRepository struct {
	db *sqlx.DB
}

func NewFatRepository(db *sqlx.DB) *fatRepository {
	return &fatRepository{db}
}

func (repo *fatRepository) AllFat(ctx context.Context, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	err := repo.db.SelectContext(ctx, &res, constants.SQL_SELECT_ALL_FATS, limit, offset)
	return res, err
}

func (repo *fatRepository) AddFat(ctx context.Context, tao *entity.Fat) error {
	_, err := repo.db.NamedExecContext(ctx, constants.SQL_INSERT_FAT, tao)
	return err
}

func (repo *fatRepository) DeleteOne(ctx context.Context, id int32) error {
	_, err := repo.db.ExecContext(ctx, constants.SQL_DELETE_FAT_BY_ID, id)
	return err
}

func (repo *fatRepository) GetByID(ctx context.Context, id int32) (entity.Fat, error) {
	var res entity.Fat
	err := repo.db.GetContext(ctx, &res, constants.SQL_SELECT_FAT_BY_ID, id)
	return res, err
}

func (repo *fatRepository) GetStates(ctx context.Context) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT state FROM fats ORDER BY state;`
	err := repo.db.SelectContext(ctx, &res, query)
	return res, err
}

func (repo *fatRepository) GetMunicipality(ctx context.Context, state string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT municipality FROM fats WHERE state = $1 ORDER BY municipality;`
	err := repo.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (repo *fatRepository) GetCounty(ctx context.Context, state, municipality string) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT county FROM fats WHERE state = $1 AND municipality = $2 ORDER BY county;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality)
	return res, err
}

func (repo *fatRepository) GetFatsByStates(ctx context.Context, state string) ([]entity.Fat, error) {
	var res []entity.Fat
	query := `SELECT * FROM fats WHERE state = $1 ORDER BY id;`
	err := repo.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (repo *fatRepository) GetFatsByMunicipality(ctx context.Context, state, municipality string) ([]entity.Fat, error) {
	var res []entity.Fat
	query := `SELECT * FROM fats WHERE state = $1 AND municipality = $2 ORDER BY id;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality)
	return res, err
}

func (repo *fatRepository) GetFatsByCounty(ctx context.Context, state, municipality, county string) ([]entity.Fat, error) {
	var res []entity.Fat
	query := `SELECT * FROM fats WHERE state = $1 AND municipality = $2 AND county = $3 ORDER BY id;`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality, county)
	return res, err
}
