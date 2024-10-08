package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
)

type FatRepository interface {
	Add(ctx context.Context, fat *entity.Fat) error
	Get(ctx context.Context, id uint) (*entity.Fat, error)
	Update(ctx context.Context, fat *entity.Fat) error
}
