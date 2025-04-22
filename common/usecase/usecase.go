package usecase

import (
	"github.com/metalpoch/olt-blueprint/common/model"
)

type InterfaceUsecase interface {
	Upsert(element *model.Interface) error
	GetAll() ([]*model.Interface, error)
	GetAllByDevice(id uint) ([]*model.Interface, error)
}

type TrafficUsecase interface {
	Add(measurement *model.Traffic) error
}

type LocationUsecase interface {
	Add(location *model.Location) (uint, error)
	FindID(location *model.Location) (uint, error)
}
