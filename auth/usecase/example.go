package usecase

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/auth/entity"
	"github.com/metalpoch/olt-blueprint/auth/model"
	"github.com/metalpoch/olt-blueprint/auth/repository"
)

type exampleUsecase struct {
	repo repository.ExampleRepository
}

func NewExampleUsecase(repo repository.ExampleRepository) *exampleUsecase {
	return &exampleUsecase{repo}
}

func (use exampleUsecase) Create(newExample *model.NewExample) (*model.Example, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	data := entity.Example{
		Message: newExample.Message,
		Value:   newExample.Value,
	}

	res, err := use.repo.Create(ctx, &data)
	if err != nil {
		return nil, err
	}

	return &model.Example{Message: res}, nil
}

func (use exampleUsecase) Get(id uint8) (*model.Example, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &model.Example{Message: res.Message}, nil
}
