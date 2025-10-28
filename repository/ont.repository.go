package repository

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

var ErrCannotDeleteEnabledONT = errors.New("cannot delete enabled ONT")

type OntRepository interface {
	Create(ctx context.Context, ont entity.Ont) error
	GetByID(ctx context.Context, id int32) (entity.Ont, error)
	GetAll(ctx context.Context) ([]entity.Ont, error)
	Delete(ctx context.Context, id int32) error
	Enable(ctx context.Context, id int32) error
	Disable(ctx context.Context, id int32) error
}

type ontRepository struct {
	db *sqlx.DB
}

func NewOntRepository(db *sqlx.DB) *ontRepository {
	return &ontRepository{db}
}

func (repo *ontRepository) Create(ctx context.Context, ont entity.Ont) error {
	query := `INSERT INTO onts (ip, ont_idx, serial, despt, line_prof, description, status, last_check)
	VALUES (:ip, :ont_idx, :serial, :despt, :line_prof, :description, :status, :last_check)`
	_, err := repo.db.NamedExecContext(ctx, query, ont)
	return err
}

func (repo *ontRepository) GetByID(ctx context.Context, id int32) (entity.Ont, error) {
	var ont entity.Ont
	query := `SELECT * FROM onts WHERE id = $1`
	err := repo.db.GetContext(ctx, &ont, query, id)
	return ont, err
}

func (repo *ontRepository) GetAll(ctx context.Context) ([]entity.Ont, error) {
	var onts []entity.Ont
	query := `SELECT * FROM onts ORDER BY id`
	err := repo.db.SelectContext(ctx, &onts, query)
	return onts, err
}

func (repo *ontRepository) Delete(ctx context.Context, id int32) error {
	query := `DELETE FROM onts WHERE id = $1 AND enabled = false`
	result, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrCannotDeleteEnabledONT
	}

	return nil
}

func (repo *ontRepository) Enable(ctx context.Context, id int32) error {
	query := `UPDATE onts SET enabled = true WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *ontRepository) Disable(ctx context.Context, id int32) error {
	query := `UPDATE onts SET enabled = false WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

