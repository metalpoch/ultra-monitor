package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/constants"
)

type UserRepository interface {
	Create(ctx context.Context, user entity.User) error
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
	_, err := repo.db.NamedExecContext(ctx, constants.SQL_USER_BY_ID, user)
	return err
}

func (repo *userRepository) UserByID(ctx context.Context, id int) (entity.User, error) {
	var res entity.User
	err := repo.db.GetContext(ctx, &res, constants.SQL_USER_BY_ID, id)
	return res, err
}

func (repo *userRepository) UserByUsername(ctx context.Context, username string) (entity.User, error) {
	var res entity.User
	err := repo.db.GetContext(ctx, &res, constants.SQL_USER_BY_USERNAME, username)
	return res, err
}

func (repo *userRepository) Disable(ctx context.Context, id int) error {
	_, err := repo.db.ExecContext(ctx, constants.SQL_DISABLE_USER, id)
	return err
}

func (repo *userRepository) Enable(ctx context.Context, id int) error {
	_, err := repo.db.ExecContext(ctx, constants.SQL_ENABLE_USER, id)
	return err
}

func (repo *userRepository) ChangePassword(ctx context.Context, id int, password string) error {
	_, err := repo.db.ExecContext(ctx, constants.SQL_CHANGE_PASSWORD, password, id)
	return err
}
