package usecase

import (
	"github.com/metalpoch/olt-blueprint/measurement/model"
)

type TemplateUsecase interface {
	Add(template *model.AddTemplate) error
	GetByID(id uint) (model.Template, error)
	GetAll() ([]model.Template, error)
	Update(id uint, template *model.AddTemplate) error
	Delete(id uint) error
}

type DeviceUsecase interface {
	Add(device *model.AddDevice) error
	Check(device *model.Device) error
	GetAll() ([]*model.Device, error)
	GetDeviceWithOIDRows() ([]*model.DeviceWithOID, error)
	Update(id uint, device *model.AddDevice) error
	Delete(id uint) error
}

type InterfaceUsecase interface {
	Upsert(element *model.Interface) error
	GetAll() ([]*model.Interface, error)
	GetAllByDevice(id uint) ([]*model.Interface, error)
}

type MeasurementUsecase interface {
	Get(id uint) (*model.Measurement, error)
	Upsert(measurement *model.Measurement) error
}

type TrafficUsecase interface {
	Add(measurement *model.Traffic) error
}
