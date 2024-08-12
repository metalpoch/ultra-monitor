package usecase

import (
	"github.com/metalpoch/olt-blueprint/update/model"
)

type DeviceUsecase interface {
	Create(ip, community string) error
	FindAll() ([]*model.Device, error)
}

type OltElementUsecase interface {
	Create(element model.ElementOLT) (string, error)
}

type CountUsecase interface {
	Create(count model.Count) (string, error)
	Find(olt, interfaceName string) (*model.Count, error)
}

type TrafficUsecase interface {
	Create(count model.CountDiff) (string, error)
}
