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

type MeasurementUsecase interface {
	Get(id uint) (*model.Measurement, error)
	Upsert(measurement *model.Measurement) error
}
