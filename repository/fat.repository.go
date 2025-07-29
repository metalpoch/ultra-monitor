package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type FatRepository interface {
	AllInfo(ctx context.Context, page, limit uint16) ([]entity.Fat, error)
	AddInfo(ctx context.Context, fat entity.Fat) (int64, error)
	DeleteOne(ctx context.Context, id int32) error
	FindByID(ctx context.Context, id int32) (entity.Fat, error)
	FindByStates(ctx context.Context, state string, page, limit uint16) ([]entity.Fat, error)
	FindByMunicipality(ctx context.Context, state, municipality string, page, limit uint16) ([]entity.Fat, error)
	FindByCounty(ctx context.Context, state, municipality, county string, page, limit uint16) ([]entity.Fat, error)
	FindBytOdn(ctx context.Context, state, municipality, county, odn string, page, limit uint16) ([]entity.Fat, error)
}

type fatRepository struct {
	db *sqlx.DB
}

func NewFatRepository(db *sqlx.DB) *fatRepository {
	return &fatRepository{db}
}

func (r *fatRepository) AllInfo(ctx context.Context, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	query := `SELECT * FROM fats ORDER BY region, state, municipality, county LIMIT $1 OFFSET $2;`
	err := r.db.SelectContext(ctx, &res, query, limit, offset)
	return res, err
}

func (r *fatRepository) AddInfo(ctx context.Context, fat entity.Fat) (int64, error) {
	var id int64
	query := `
        INSERT INTO fats (ip, region, state, municipality, county, odn, fat, shell, card, port)
        VALUES (:ip, :region, :state, :municipality, :county, :odn, :fat, :shell, :card, :port)
        RETURNING id;`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &id, fat)
	return id, err
}

func (r *fatRepository) DeleteOne(ctx context.Context, id int32) error {
	query := `DELETE FROM fats WHERE id = $1;`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *fatRepository) FindByID(ctx context.Context, id int32) (entity.Fat, error) {
	var fat entity.Fat
	query := `SELECT * FROM fats WHERE id = $1;`
	err := r.db.GetContext(ctx, &fat, query, id)
	if err != nil {
		return entity.Fat{}, err
	}
	return fat, nil
}

func (r *fatRepository) FindByStates(ctx context.Context, state string, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	query := `SELECT * FROM fats WHERE state = $1 ORDER BY region, state, municipality, county LIMIT $2 OFFSET $3;`
	err := r.db.SelectContext(ctx, &res, query, state, limit, offset)
	return res, err
}

func (r *fatRepository) FindByMunicipality(ctx context.Context, state, municipality string, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	query := `SELECT * FROM fats WHERE state = $1 AND municipality = $2 ORDER BY region, state, municipality, county LIMIT $3 OFFSET $4;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, limit, offset)
	return res, err
}

func (r *fatRepository) FindByCounty(ctx context.Context, state, municipality, county string, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	query := `SELECT * FROM fats WHERE state = $1 AND municipality = $2 AND county = $3 ORDER BY region, state, municipality, county LIMIT $4 OFFSET $5;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county, limit, offset)
	return res, err
}

func (r *fatRepository) FindBytOdn(ctx context.Context, state, municipality, county, odn string, page, limit uint16) ([]entity.Fat, error) {
	var res []entity.Fat
	offset := (page - 1) * limit
	query := `SELECT * FROM fats WHERE state = $1 AND municipality = $2 AND county = $3 AND odn = $4 ORDER BY region, state, municipality, county LIMIT $5 OFFSET $6;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county, odn, limit, offset)
	return res, err
}
