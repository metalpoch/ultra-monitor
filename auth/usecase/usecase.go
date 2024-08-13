package usecase

import "github.com/metalpoch/olt-blueprint/auth/model"

type ExampleUsecase interface {
	Create(newExample *model.NewExample) (*model.Example, error)
	Get(id uint8) (*model.Example, error)
}
