package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/auth/entity"
)

type ExampleRepository interface {
	Create(ctx context.Context, data *entity.Example) (string, error)
	Get(ctx context.Context, id uint8) (*entity.ExampleResponse, error)
}
