package usecase

import (
	"github.com/metalpoch/olt-blueprint/update/model"
)

type TemplateUsecase interface {
	Add(template *model.AddTemplate) error
	GetByID(id uint) (model.Template, error)
	GetAll() ([]model.Template, error)
}

type DeviceUsecase interface {
	Add(device *model.AddDevice) error
	Check(device *model.CheckDevice) error
	GetAll() ([]*model.Device, error)
	GetDeviceWithOIDRows() ([]*model.DeviceWithOID, error)
}

type InterfaceUsecase interface {
	Upsert(element *model.Interface) error
	GetAll() ([]*model.Interface, error)
	GetAllByDevice(id uint) ([]*model.Interface, error)
}

type MeasurementUsecase interface {
	Get(id uint) error
	Upsert(measurement *model.Measurement) error
}
