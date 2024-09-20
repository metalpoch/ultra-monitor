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
	GetAll() ([]model.Device, error)
}
