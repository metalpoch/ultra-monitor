package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetAll(ctx context.Context) ([]*entity.User, error)
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	SoftDelete(ctx context.Context, id uint) error
	ChangePassword(ctx context.Context, id uint, password string) error
}
