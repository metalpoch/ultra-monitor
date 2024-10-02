package usecase

import (
	"github.com/metalpoch/olt-blueprint/traffic/model"
)

type TrafficUsecase interface {
	GetTrafficByInterface(id uint, date *model.RangeDate) ([]*model.Traffic, error)
	GetTrafficByDevice(id uint, date *model.RangeDate) ([]*model.Traffic, error)
}
