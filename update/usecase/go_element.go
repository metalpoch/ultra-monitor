package usecase

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/repository"
)

type elementUsecase struct {
	repo repository.ElementRepository
}

func NewElementsUsecase(repo repository.ElementRepository) *elementUsecase {
	return &elementUsecase{repo}
}

func (use elementUsecase) Create(element model.ElementOLT) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	id, err := use.repo.Create(ctx, &entity.ElementOLT{
		OLT:       element.OLT,
		Interface: element.Interface,
		Slot:      element.Slot,
		Card:      element.Card,
		Port:      element.Port,
	})

	if err != nil {
		return "", err
	}

	return id, nil
}

func (use elementUsecase) FindID(element model.ElementOLT) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return use.repo.FindID(ctx, element.OLT, element.Interface)
}
