package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type OltRepository interface {
	Add(ctx context.Context, olt *entity.Olt) error
	Olt(ctx context.Context, id string) (entity.Olt, error)
	Olts(ctx context.Context) ([]entity.Olt, error)
	Delete(ctx context.Context, id string) error
	GetAllIP(ctx context.Context) ([]string, error)
	GetAllSysname(ctx context.Context) ([]string, error)
	OltsByState(ctx context.Context, state string) ([]entity.Olt, error)
	OltsByMunicipality(ctx context.Context, state, municipality string) ([]entity.Olt, error)
	OltsByCounty(ctx context.Context, state, municipality, county string) ([]entity.Olt, error)
}

type oltRepository struct {
	db *sqlx.DB
}

func NewOltRepository(db *sqlx.DB) *oltRepository {
	return &oltRepository{db}
}

func (repo *oltRepository) Add(ctx context.Context, olt *entity.Olt) error {
	query := `
	INSERT INTO olts (ip, community, sys_name, sys_location, is_alive, last_check)
	VALUES (:ip, :community, :sys_name, :sys_location, :is_alive, :last_check)`
	_, err := repo.db.NamedExecContext(ctx, query, olt)
	return err
}

func (repo *oltRepository) Olt(ctx context.Context, id string) (entity.Olt, error) {
	var olt entity.Olt
	query := `SELECT * FROM olts WHERE ip = $1`
	err := repo.db.GetContext(ctx, &olt, query, id)
	if err != nil {
		return olt, err
	}
	return olt, nil
}
func (repo *oltRepository) Olts(ctx context.Context) ([]entity.Olt, error) {
	var res []entity.Olt
	query := `SELECT * FROM olts ORDER BY sys_name`
	err := repo.db.GetContext(ctx, &res, query)
	return res, err
}

func (repo *oltRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM olts WHERE ip = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *oltRepository) OltsByState(ctx context.Context, state string) ([]entity.Olt, error) {
	var res []entity.Olt
	query := `
    SELECT DISTINCT olts.*
    FROM olts
    JOIN fats ON fats.olt_ip = olts.ip
    WHERE fats.state = $1
    ORDER BY olts.sys_name`
	err := repo.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (repo *oltRepository) OltsByMunicipality(ctx context.Context, state, municipality string) ([]entity.Olt, error) {
	var res []entity.Olt
	query := `
	SELECT DISTINCT olts.*
	FROM olts
	JOIN fats ON fats.olt_ip = olts.ip
	WHERE fats.state = $1 AND fats.municipality = $2
	ORDER BY olts.sys_name`
	err := repo.db.SelectContext(ctx, &res, query, state, municipality)
	return res, err
}

func (repo *oltRepository) OltsByCounty(ctx context.Context, state, municipality, county string) ([]entity.Olt, error) {
	var res []entity.Olt
	query := `
	SELECT DISTINCT olts.*
	FROM olts
	JOIN fats ON fats.olt_ip = olts.ip
	WHERE fats.state = $1 AND fats.municipality = $2 AND fats.county = $3
	ORDER BY olts.sys_name`

	err := repo.db.SelectContext(ctx, &res, query, state, municipality, county)
	return res, err
}

func (repo *oltRepository) GetAllIP(ctx context.Context) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT ip FROM olts ORDER BY ip`
	err := repo.db.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (repo *oltRepository) GetAllSysname(ctx context.Context) ([]string, error) {
	var res []string
	query := `SELECT DISTINCT sysname FROM olts ORDER BY sysname`
	err := repo.db.SelectContext(ctx, &res, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}
