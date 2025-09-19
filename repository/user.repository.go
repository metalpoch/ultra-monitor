package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) error
	AllUsers(ctx context.Context) ([]entity.UserResponse, error)
	UserByID(ctx context.Context, id int) (entity.User, error)
	UserByUsername(ctx context.Context, username string) (entity.User, error)
	Disable(ctx context.Context, id int) error
	Enable(ctx context.Context, id int) error
	ChangePassword(ctx context.Context, id int, password string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db}
}

func (repo *userRepository) Create(ctx context.Context, user entity.User) error {
	query := `INSERT INTO users (id, fullname, username, password, change_password, is_admin, is_disabled)
	VALUES (:id, :fullname, :username, :password, :change_password, :is_admin, :is_disabled);`
	_, err := repo.db.NamedExecContext(ctx, query, user)
	return err
}

func (repo *userRepository) AllUsers(ctx context.Context) ([]entity.UserResponse, error) {
	var res []entity.UserResponse
	query := `SELECT id, fullname, username, is_admin, is_disabled, change_password, created_at FROM users ORDER BY id`
	err := repo.db.SelectContext(ctx, &res, query)
	return res, err
}

func (repo *userRepository) UserByID(ctx context.Context, id int) (entity.User, error) {
	var res entity.User
	query := `SELECT * FROM users WHERE id = $1;`
	err := repo.db.GetContext(ctx, &res, query, id)
	return res, err
}

func (repo *userRepository) UserByUsername(ctx context.Context, username string) (entity.User, error) {
	var res entity.User
	query := `SELECT * FROM users WHERE username = $1;`
	err := repo.db.GetContext(ctx, &res, query, username)
	return res, err
}

func (repo *userRepository) Disable(ctx context.Context, id int) error {
	query := `UPDATE users SET is_disabled = true WHERE id = $1;`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *userRepository) Enable(ctx context.Context, id int) error {
	query := `UPDATE users SET is_disabled = false WHERE id = $1;`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *userRepository) ChangePassword(ctx context.Context, id int, password string) error {
	query := `UPDATE users SET password = $1, change_password = false WHERE id = $2;`
	_, err := repo.db.ExecContext(ctx, query, password, id)
	return err
}
