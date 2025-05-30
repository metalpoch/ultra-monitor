package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
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
}

type oltRepository struct {
	db *sqlx.DB
}

func NewOltRepository(db *sqlx.DB) *oltRepository {
	return &oltRepository{db}
}

func (repo *oltRepository) Add(ctx context.Context, device *entity.Olt) error {
	query := `
        INSERT INTO olts (ip, community, sys_name, sys_location, is_alive, last_check)
        VALUES (:ip, :community, :sys_name, :sys_location, :is_alive, :last_check)
    `
	_, err := repo.db.NamedExecContext(ctx, query, device)
	return err
}

func (repo *oltRepository) Olt(ctx context.Context, id int32) (entity.Olt, error) {
	var olt entity.Olt
	query := `SELECT * FROM olts WHERE id = $1`
	err := repo.db.GetContext(ctx, &olt, query, id)
	if err != nil {
		return olt, err
	}
	return olt, nil
}

func (repo *oltRepository) Update(ctx context.Context, olt entity.Olt) error {
	query := `
        UPDATE olts SET
            ip = :ip,
            community = :community,
            sys_name = :sys_name,
            sys_location = :sys_location,
            is_alive = :is_alive,
            last_check = :last_check
        WHERE id = :id
    `
	_, err := repo.db.NamedExecContext(ctx, query, olt)
	return err
}

func (repo *oltRepository) Delete(ctx context.Context, id int32) error {
	query := `DELETE FROM olts WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *oltRepository) Olts(ctx context.Context, page, limit uint16) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	query := `SELECT * FROM olts ORDER BY sys_name LIMIT $1 OFFSET $2`
	err := repo.db.SelectContext(ctx, &res, query, limit, offset)
	return res, err
}

func (repo *oltRepository) OltsByState(ctx context.Context, state string, page, limit uint16) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	query := `
		SELECT DISTINCT olts.*
		FROM olts
		JOIN fats ON fats.olt_ip = olts.ip
		WHERE fats.state = $1
		ORDER BY olts.sys_name
		LIMIT $2 OFFSET $3`
	err := repo.db.SelectContext(ctx, &res, query, state, limit, offset)
	return res, err
}

func (repo *oltRepository) OltsByCounty(ctx context.Context, state, county string, page, limit uint16) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	query := `
		SELECT DISTINCT olts.*
		FROM olts
		JOIN fats ON fats.olt_ip = olts.ip
		WHERE fats.state = $1 AND fats.county = $2
		ORDER BY olts.sys_name
		LIMIT $3 OFFSET $4`
	err := repo.db.SelectContext(ctx, &res, query, state, county, limit, offset)
	return res, err
}

func (repo *oltRepository) OltsByMunicipality(ctx context.Context, state, county, municipality string, page, limit uint16) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	query := `
		SELECT DISTINCT olts.*
		FROM olts
		JOIN fats ON fats.olt_ip = olts.ip
		WHERE fats.state = $1 AND fats.county = $2 AND fats.municipality = $3
		ORDER BY olts.sys_name
		LIMIT $4 OFFSET $5`
	err := repo.db.SelectContext(ctx, &res, query, state, county, municipality, limit, offset)
	return res, err
}
