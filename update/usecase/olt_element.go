package usecase

import (
	"context"
	"time"

	"github.com/metalpoch/olt-blueprint/update/entity"
	"github.com/metalpoch/olt-blueprint/update/model"
	"github.com/metalpoch/olt-blueprint/update/repository"
)

type oltElementUsecase struct {
	repo repository.ElementRepository
}

func NewOltElementUsecase(repo repository.ElementRepository) *oltElementUsecase {
	return &oltElementUsecase{repo}
}

func (use oltElementUsecase) Create(element model.ElementOLT) (string, error) {
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

// func (use oltElementUsecase) FindID(element model.Element) (string, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	id, err := use.repo.FindID(ctx, entity.Element{
// 		Sysname: element.Sysname,
// 		Slot:    element.Slot,
// 		Card:    element.Card,
// 		Port:    element.Port,
// 	})

// 	if err != nil {
// 		return "", err
// 	}

// 	return id, nil
// }

// func (use oltElementUsecase) FindByID(id string) (*model.Element, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	res, err := use.repo.FindByID(ctx, id)
// 	if err != nil {
// 		return nil, err
// 	}

// 	element := &model.Element{
// 		ID:      res.ID.Hex(),
// 		Sysname: res.Sysname,
// 		Slot:    res.Slot,
// 		Card:    res.Card,
// 		Port:    res.Port,
// 	}

// 	return element, nil
// }

// func (use oltElementUsecase) FindBySysname(sysname string) ([]*model.Element, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	res, err := use.repo.FindBySysname(ctx, sysname)
// 	if err != nil {
// 		return nil, err
// 	}

// 	elements := []*model.Element{}
// 	for _, e := range res {
// 		elements = append(elements, &model.Element{
// 			ID:      e.ID.Hex(),
// 			Sysname: e.Sysname,
// 			Slot:    e.Slot,
// 			Card:    e.Card,
// 			Port:    e.Port,
// 		})
// 	}

// 	return elements, nil
// }
