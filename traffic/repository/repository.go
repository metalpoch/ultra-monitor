package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type FeedRepository interface {
	GetDevice(ctx context.Context, id uint) (*entity.Device, error)
	GetInterface(ctx context.Context, id uint) (*entity.Interface, error)
	GetAllDevice(ctx context.Context) ([]*entity.Device, error)
}
