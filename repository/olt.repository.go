package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/constants"
)

type OltRepository interface {
	Add(ctx context.Context, olt *entity.Olt) error
	Update(ctx context.Context, olt entity.Olt) error
	Delete(ctx context.Context, id int32) error
	Olt(ctx context.Context, id int32) (entity.Olt, error)
	Olts(ctx context.Context, page, limit uint16) ([]entity.Olt, error)
	OltsByState(ctx context.Context, state string, page, limit uint16) ([]entity.Olt, error)
	OltsByCounty(ctx context.Context, state, county string, page, limit uint16) ([]entity.Olt, error)
	OltsByMunicipality(ctx context.Context, state, county, municipality string, page, limit uint16) ([]entity.Olt, error)
	GetIPs(ctx context.Context) ([]string, error)
}

type oltRepository struct {
	db *sqlx.DB
}

func NewOltRepository(db *sqlx.DB) *oltRepository {
	return &oltRepository{db}
}

func (repo *oltRepository) Add(ctx context.Context, olt *entity.Olt) error {
	_, err := repo.db.NamedExecContext(ctx, constants.SQL_ADD_OLT, olt)
	return err
}

func (repo *oltRepository) Olt(ctx context.Context, id int32) (entity.Olt, error) {
	var olt entity.Olt
	err := repo.db.GetContext(ctx, &olt, constants.SQL_GET_OLT, id)
	if err != nil {
		return olt, err
	}
	return olt, nil
}

func (repo *oltRepository) Update(ctx context.Context, olt entity.Olt) error {
	_, err := repo.db.NamedExecContext(ctx, constants.SQL_UPDATE_OLT, olt)
	return err
}

func (repo *oltRepository) Delete(ctx context.Context, id int32) error {
	_, err := repo.db.ExecContext(ctx, constants.SQL_DELETE_OLT, id)
	return err
}

func (repo *oltRepository) Olts(ctx context.Context, page, limit uint16) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	err := repo.db.SelectContext(ctx, &res, constants.SQL_GET_ALL_OLTS, limit, offset)
	return res, err
}

func (repo *oltRepository) OltsByState(ctx context.Context, state string, page, limit uint16) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	err := repo.db.SelectContext(ctx, &res, constants.SQL_GET_OLTS_BY_STATE, state, limit, offset)
	return res, err
}

func (repo *oltRepository) OltsByCounty(ctx context.Context, state, county string, page, limit uint16) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	err := repo.db.SelectContext(ctx, &res, constants.SQL_GET_OLTS_BY_COUNTY, state, county, limit, offset)
	return res, err
}

func (repo *oltRepository) OltsByMunicipality(ctx context.Context, state, county, municipality string, page, limit uint16) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	err := repo.db.SelectContext(ctx, &res, constants.SQL_GET_OLTS_BY_MUNICIPALITY, state, county, municipality, limit, offset)
	return res, err
}

func (repo *oltRepository) GetIPs(ctx context.Context) ([]string, error) {
	var res []string
	err := repo.db.SelectContext(ctx, &res, constants.SQL_GET_OLTS_IPS)
	if err != nil {
		return nil, err
	}
	return res, nil
}
