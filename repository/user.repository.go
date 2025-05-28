package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) error
	UserByID(ctx context.Context, id uint32) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	Disable(ctx context.Context, id uint32) error
	Enable(ctx context.Context, id uint32) error
	ChangePassword(ctx context.Context, id uint32, password string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *userRepository {
	return &userRepository{db}
}

func (repo userRepository) Create(ctx context.Context, user entity.User) error {
	query := `
	INSERT INTO users (id, fullname, email, password, change_password, is_admin, is_disabled, created_at)
	VALUES (:id, :fullname, :email, :password, :change_password, :is_admin, :is_disabled, :created_at)`
	_, err := repo.db.NamedExecContext(ctx, query, user)
	return err
}

func (repo userRepository) UserByID(ctx context.Context, id uint32) (entity.User, error) {
	var res entity.User
	query := `SELECT id, fullname, email, password, change_password, is_admin, is_disabled, created_at FROM users WHERE id = $1`
	err := repo.db.GetContext(ctx, &res, query, id)
	return res, err
}

func (repo userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var res entity.User
	query := `SELECT id, fullname, email, password, change_password, is_admin, is_disabled, created_at FROM users WHERE email = $1`
	err := repo.db.GetContext(ctx, &res, query, email)
	return res, err
}

func (repo userRepository) Disable(ctx context.Context, id uint32) error {
	query := `UPDATE users SET is_disabled = 0 WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo userRepository) Enable(ctx context.Context, id uint32) error {
	query := `UPDATE users SET is_disabled = 1 WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo userRepository) ChangePassword(ctx context.Context, id uint32, password string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`
	_, err := repo.db.ExecContext(ctx, query, password, id)
	return err
}
