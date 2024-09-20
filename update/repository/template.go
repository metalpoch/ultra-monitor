package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
)

type templateRepository struct {
	db *sql.DB
}

func NewTemplateRepository(db *sql.DB) *templateRepository {
	return &templateRepository{db}
}

func (repo templateRepository) Add(ctx context.Context, template *entity.Template) error {
	q := `INSERT INTO template
		(name, oid_bw, oid_in, oid_out, prefix_bw, prefix_in, prefix_out, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING id;`
	now := time.Now()
	if err := repo.db.QueryRowContext(
		ctx,
		q,
		template.Name,
		template.OidBw,
		template.OidIn,
		template.OidOut,
		template.PrefixBw,
		template.PrefixIn,
		template.PrefixOut,
		now,
		now,
	).Scan(&template.ID); err != nil {
		return err
	}

	return nil
}

func (repo templateRepository) GetByID(ctx context.Context, id uint) (*entity.Template, error) {
	t := new(entity.Template)
	row := repo.db.QueryRowContext(ctx, "SELECT * FROM template WHERE id=$1;", id)

	err := row.Scan(
		&t.ID,
		&t.Name,
		&t.OidBw,
		&t.OidIn,
		&t.OidOut,
		&t.PrefixBw,
		&t.PrefixIn,
		&t.PrefixOut,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return t, nil

}
func (repo templateRepository) GetAll(ctx context.Context) ([]*entity.Template, error) {
	templates := []*entity.Template{}
	rows, err := repo.db.QueryContext(ctx, "SELECT * FROM template;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		t := new(entity.Template)
		rows.Scan(
			&t.ID,
			&t.Name,
			&t.OidBw,
			&t.OidIn,
			&t.OidOut,
			&t.PrefixBw,
			&t.PrefixIn,
			&t.PrefixOut,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		templates = append(templates, t)
	}

	return templates, nil
}
