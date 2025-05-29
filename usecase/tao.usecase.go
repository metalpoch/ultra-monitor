package usecase

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/entity"
	"github.com/metalpoch/olt-blueprint/internal/dto"
	"github.com/metalpoch/olt-blueprint/model"
	"github.com/metalpoch/olt-blueprint/repository"
)

type TaoUsecase struct {
	repo repository.TaoRepository
}

func NewTaoUsecase(db *sqlx.DB) *TaoUsecase {
	return &TaoUsecase{repository.NewTaoRepository(db)}
}

func (uc *TaoUsecase) AddTao(tao model.Tao) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	entityTao := (entity.Tao)(tao)
	err := uc.repo.AddTao(ctx, &entityTao)
	if err != nil {
		return err
	}

	return nil
}

func (uc *TaoUsecase) DeleteOne(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := uc.repo.DeleteOne(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *TaoUsecase) GetAll(pag dto.Pagination) ([]model.Tao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.AllTao(ctx, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var taos []model.Tao
	for _, t := range res {
		taos = append(taos, (model.Tao)(t))
	}

	return taos, nil
}

func (uc *TaoUsecase) GetByID(id uint64) (model.Tao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return model.Tao{}, err
	}

	return (model.Tao)(res), nil
}
func (uc *TaoUsecase) GetByTao(tao string) (model.Tao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetByTao(ctx, tao)
	if err != nil {
		return model.Tao{}, err
	}

	return (model.Tao)(res), nil
}
func (uc *TaoUsecase) GetByOdn(state, odn string) ([]model.Tao, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetByOdn(ctx, state, odn)
	if err != nil {
		return nil, err
	}

	var taos []model.Tao
	for _, t := range res {
		taos = append(taos, (model.Tao)(t))
	}

	return taos, nil
}

func (uc *TaoUsecase) GetOdnStates(state string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetOdnStates(ctx, state)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (uc *TaoUsecase) GetOdnCounty(state, county string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetOdnCounty(ctx, state, county)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (uc *TaoUsecase) GetOdnMunicipality(state, county, municipality string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetOdnMunicipality(ctx, state, county, municipality)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (uc *TaoUsecase) GetOdnByOlt(oltIP string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetOdnByOlt(ctx, oltIP)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (uc *TaoUsecase) GetOdnOltPort(oltIP string, shell, card, port uint8) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.GetOdnOltPort(ctx, oltIP, shell, card, port)
	if err != nil {
		return nil, err
	}

	return res, nil
}
