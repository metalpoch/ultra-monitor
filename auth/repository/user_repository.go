package repository

import (
	"context"
	"database/sql"

	"github.com/metalpoch/olt-blueprint/auth/entity"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db}
}

func (repo userRepository) Create(ctx context.Context, user *entity.User) error {
	q := `
    INSERT INTO users (
		id,
		first_name,
		last_name,
		email,
		password,
		change_password,
		states_permission,
		is_admin)
    	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    	RETURNING id;
    `
	if err := repo.db.QueryRowContext(
		ctx, q, user.ID, user.FirstName, user.LastName, user.Email, user.Password,
		user.ChangePassword, user.StatesPermission, user.IsAdmin,
	).Scan(&user.ID); err != nil {
		return err
	}

	return nil
}

func (repo userRepository) GetAll(ctx context.Context) ([]*entity.User, error) {
	users := []*entity.User{}
	q := "SELECT * FROM users WHERE is_disable=false;"

	rows, err := repo.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		u := new(entity.User)
		rows.Scan(
			u.ID,
			u.FirstName,
			u.LastName,
			u.Email,
			u.Password,
			u.ChangePassword,
			u.StatesPermission,
			u.IsAdmin,
			u.IsDisabled,
			u.CreatedAt,
			u.UpdatedAt,
		)
		users = append(users, u)
	}

	return users, nil
}

func (repo userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := new(entity.User)
	row := repo.db.QueryRowContext(ctx, "SELECT * FROM users WHERE email=$1;", email)

	err := row.Scan(
		user.ID,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.ChangePassword,
		user.StatesPermission,
		user.IsAdmin,
		user.IsDisabled,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repo userRepository) SoftDelete(ctx context.Context, id uint) error {
	q := "UPDATE users set is_disabled=true  WHERE id=$1;"

	stmt, err := repo.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, id); err != nil {
		return err
	}

	return nil
}

func (repo userRepository) ChangePassword(ctx context.Context, id, password string) error {
	q := "UPDATE users set change_password=false, password=$1 WHERE id=$2 AND change_password=true;"

	stmt, err := repo.db.PrepareContext(ctx, q)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.ExecContext(ctx, password, id); err != nil {
		return err
	}

	return nil
}
