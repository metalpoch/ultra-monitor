package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type InfoDeviceRepository interface {
	AllInfo(ctx context.Context, page, limit uint16) ([]entity.InfoDevice, error)
	AddInfo(ctx context.Context, fat entity.InfoDevice) (int64, error)
	DeleteOne(ctx context.Context, id int32) error
	FindByID(ctx context.Context, id int32) (entity.InfoDevice, error)
	FindByStates(ctx context.Context, state string) ([]entity.InfoDevice, error)
	FindByMunicipality(ctx context.Context, state, municipality string) ([]entity.InfoDevice, error)
	FindByCounty(ctx context.Context, state, municipality, county string) ([]entity.InfoDevice, error)
	FindBytOdn(ctx context.Context, state, municipality, county, odn string) ([]entity.InfoDevice, error)
}

type infoDeviceRepository struct {
	db *sqlx.DB
}

func NewInfoDeviceRepository(db *sqlx.DB) *infoDeviceRepository {
	return &infoDeviceRepository{db}
}

func (r *infoDeviceRepository) AllInfo(ctx context.Context, page, limit uint16) ([]entity.InfoDevice, error) {
	var res []entity.InfoDevice
	offset := (page - 1) * limit
	query := `SELECT * FROM info_device ORDER BY region, state, municipality, county LIMIT $1 OFFSET $2;`
	err := r.db.SelectContext(ctx, &res, query, limit, offset)
	return res, err
}

func (r *infoDeviceRepository) AddInfo(ctx context.Context, fat entity.InfoDevice) (int64, error) {
	var id int64
	query := `
        INSERT INTO info_device (ip, region, state, municipality, county, odn, fat, pon_shell, pon_card, pon_port)
        VALUES (:ip, :region, :state, :municipality, :county, :odn, :fat, :pon_shell, :pon_card, :pon_port)
        RETURNING id;`
	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &id, fat)
	return id, err
}

func (r *infoDeviceRepository) DeleteOne(ctx context.Context, id int32) error {
	query := `DELETE FROM info_device WHERE id = $1;`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *infoDeviceRepository) FindByID(ctx context.Context, id int32) (entity.InfoDevice, error) {
	var fat entity.InfoDevice
	query := `SELECT * FROM info_device WHERE id = $1;`
	err := r.db.GetContext(ctx, &fat, query, id)
	if err != nil {
		return entity.InfoDevice{}, err
	}
	return fat, nil
}

func (r *infoDeviceRepository) FindByStates(ctx context.Context, state string) ([]entity.InfoDevice, error) {
	var res []entity.InfoDevice
	query := `SELECT * FROM info_device WHERE state = $1 ORDER BY region, municipality, county;`
	err := r.db.SelectContext(ctx, &res, query, state)
	return res, err
}

func (r *infoDeviceRepository) FindByMunicipality(ctx context.Context, state, municipality string) ([]entity.InfoDevice, error) {
	var res []entity.InfoDevice
	query := `SELECT * FROM info_device WHERE state = $1 AND municipality = $2 ORDER BY region, county;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality)
	return res, err
}

func (r *infoDeviceRepository) FindByCounty(ctx context.Context, state, municipality, county string) ([]entity.InfoDevice, error) {
	var res []entity.InfoDevice
	query := `SELECT * FROM info_device WHERE state = $1 AND municipality = $2 AND county = $3 ORDER BY region;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county)
	return res, err
}

func (r *infoDeviceRepository) FindBytOdn(ctx context.Context, state, municipality, county, odn string) ([]entity.InfoDevice, error) {
	var res []entity.InfoDevice
	query := `SELECT * FROM info_device WHERE state = $1 AND municipality = $2 AND county = $3 AND odn = $4 ORDER BY region;`
	err := r.db.SelectContext(ctx, &res, query, state, municipality, county, odn)
	return res, err
}
