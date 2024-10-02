package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/traffic/model"
)

type TrafficRepository interface {
	GetTrafficByInterface(ctx context.Context, id uint, date *model.RangeDate) ([]*entity.Traffic, error)
	GetTrafficByDevice(ctx context.Context, id uint, date *model.RangeDate) ([]*entity.Traffic, error)
}
