package repository

import (
	"context"

	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
)

type FeedRepository interface {
	GetDevice(ctx context.Context, id uint) (*entity.Device, error)
	GetAllDevice(ctx context.Context) ([]*entity.Device, error)
	GetInterface(ctx context.Context, id uint) (*entity.Interface, error)
	GetInterfacesByDevice(ctx context.Context, id uint) ([]*entity.Interface, error)
	GetLocationStates(ctx context.Context) ([]*string, error)
	GetLocationCounties(ctx context.Context, state string) ([]*string, error)
	GetLocationMunicipalities(ctx context.Context, state, county string) ([]*string, error)
}

type TrafficRepository interface {
	GetTrafficByInterface(ctx context.Context, id uint, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByDevice(ctx context.Context, id uint, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByFat(ctx context.Context, id uint, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByLocationID(ctx context.Context, id uint, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByState(ctx context.Context, state string, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByCounty(ctx context.Context, state, county string, date *model.TranficRangeDate) ([]*entity.Traffic, error)
	GetTrafficByMunicipality(ctx context.Context, state, county, municipality string, date *model.TranficRangeDate) ([]*entity.Traffic, error)
}
