package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/auth/entity"
)

type UserRepository interface {
	Create(ctx context.Context, data *entity.User) (string, error)
	Get(ctx context.Context) (*entity.Users, error)
	GetValue(ctx context.Context, clave string, valor string) (*entity.Users, error)
	DeleteName(ctx context.Context, name string) (string, error)
	ChangePassword(ctx context.Context, user *entity.User) (string, error)
	Login(ctx context.Context, email string, pasword string) (string, *entity.User, error)
}
