package usecase

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/metalpoch/ultra-monitor/repository"
)

type FatUsecase struct {
	repo repository.FatRepository
}

func NewFatUsecase(db *sqlx.DB) *FatUsecase {
	return &FatUsecase{repository.NewFatRepository(db)}
}

func (use *FatUsecase) AddFat(tao model.Fat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	entityTao := (entity.Fat)(tao)
	err := use.repo.AddFat(ctx, &entityTao)
	if err != nil {
		return err
	}

	return nil
}

func (use *FatUsecase) DeleteOne(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := use.repo.DeleteOne(ctx, int32(id))
	if err != nil {
		return err
	}

	return nil
}

func (use *FatUsecase) GetAll(pag dto.Pagination) ([]model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.AllFat(ctx, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var taos []model.Fat
	for _, t := range res {
		taos = append(taos, (model.Fat)(t))
	}

	return taos, nil
}

func (use *FatUsecase) GetByID(id int) (model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetByID(ctx, int32(id))
	if err != nil {
		return model.Fat{}, err
	}

	return (model.Fat)(res), nil
}
func (use *FatUsecase) GetByFat(tao string) (model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetByFat(ctx, tao)
	if err != nil {
		return model.Fat{}, err
	}

	return (model.Fat)(res), nil
}
func (use *FatUsecase) GetByOdn(state, odn string) ([]model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetByOdn(ctx, state, odn)
	if err != nil {
		return nil, err
	}

	var taos []model.Fat
	for _, t := range res {
		taos = append(taos, (model.Fat)(t))
	}

	return taos, nil
}

func (use *FatUsecase) GetOdnStates(state string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetOdnByStates(ctx, state)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (use *FatUsecase) GetOdnCounty(state, county string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetOdnByCounty(ctx, state, county)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (use *FatUsecase) GetOdnMunicipality(state, county, municipality string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetOdnMunicipality(ctx, state, county, municipality)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (use *FatUsecase) GetOdnByOlt(oltIP string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetOdnByOlt(ctx, oltIP)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func (use *FatUsecase) GetOdnOltPort(oltIP string, shell, card, port int) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := use.repo.GetOdnByOltPort(ctx, oltIP, uint8(shell), uint8(card), uint8(port))
	if err != nil {
		return nil, err
	}

	return res, nil
}
