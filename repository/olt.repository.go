package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/entity"
)

type OltRepository interface {
	Add(ctx context.Context, olt *entity.Olt) error
	Update(ctx context.Context, olt entity.Olt) error
	Delete(ctx context.Context, id uint64) error
	Olt(ctx context.Context, id uint64) (entity.Olt, error)
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
        INSERT INTO olt (ip, community, sys_name, sys_location, is_alive, last_check)
        VALUES (:ip, :community, :sys_name, :sys_location, :is_alive, :last_check)
    `
	_, err := repo.db.NamedExecContext(ctx, query, device)
	return err
}

func (repo *oltRepository) Olt(ctx context.Context, id uint64) (entity.Olt, error) {
	var olt entity.Olt
	query := `SELECT * FROM olt WHERE id = $1`
	err := repo.db.GetContext(ctx, &olt, query, id)
	if err != nil {
		return olt, err
	}
	return olt, nil
}

func (repo *oltRepository) Update(ctx context.Context, olt entity.Olt) error {
	query := `
        UPDATE olt SET
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

func (repo *oltRepository) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM olt WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *oltRepository) Olts(ctx context.Context, page, limit uint16) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	query := `SELECT * FROM olt ORDER BY sys_name LIMIT ? OFFSET ?`
	err := repo.db.SelectContext(ctx, &res, query, limit, offset)
	return res, err
}

func (repo *oltRepository) OltsByState(ctx context.Context, state string, page, limit uint16) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	query := `
		SELECT id, ip, community, sys_name, sys_location, is_alive, last_check, created_at
		FROM olt
		JOIN pons ON olt.id = pons.olt_id
		JOIN fats_pon ON p.id = fats_pon.pon_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON fats.location_id = locations.id
		WHERE locations.state = ?
		LIMIT ?	OFFSET ?`
	err := repo.db.SelectContext(ctx, &res, query, state, limit, offset)
	return res, err
}

func (repo *oltRepository) OltsByCounty(ctx context.Context, state, county string, page, limit uint16) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	query := `
		SELECT id, ip, community, sys_name, sys_location, is_alive, last_check, created_at
		FROM olt
		JOIN pons ON olt.id = pons.olt_id
		JOIN fats_pon ON p.id = fats_pon.pon_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON fats.location_id = locations.id
		WHERE locations.state = ? AND locations.county = ?
		LIMIT ?	OFFSET ?`
	err := repo.db.SelectContext(ctx, &res, query, state, county, limit, offset)
	return res, err
}

func (repo *oltRepository) OltsByMunicipality(ctx context.Context, state, county, municipality string, page, limit uint16) ([]entity.Olt, error) {
	var res []entity.Olt
	offset := (page - 1) * limit
	query := `
		SELECT id, ip, community, sys_name, sys_location, is_alive, last_check, created_at
		FROM olt
		JOIN pons ON olt.id = pons.olt_id
		JOIN fats_pon ON p.id = fats_pon.pon_id
		JOIN fats ON fats.id = fats_pon.fat_id
		JOIN locations ON fats.location_id = locations.id
		WHERE locations.state = ? AND locations.county = ? AND locations.municipality = ?
		LIMIT ?	OFFSET ?`
	err := repo.db.SelectContext(ctx, &res, query, state, county, municipality, limit, offset)
	return res, err
}
