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

func (use *FatUsecase) GetAll(pag dto.Pagination) ([]model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.AllInfo(ctx, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var taos []model.Fat
	for _, t := range res {
		taos = append(taos, (model.Fat)(t))
	}

	return taos, nil
}

func (use *FatUsecase) AddInfo(info dto.Fat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := use.repo.AddInfo(ctx, entity.Fat{
		Fat:          info.Fat,
		Region:       info.Region,
		State:        info.State,
		Municipality: info.Municipality,
		County:       info.County,
		Odn:          info.Odn,
		IP:           info.IP,
		Shell:        info.Shell,
		Card:         info.Card,
		Port:         info.Port,
	})

	return err
}

func (use *FatUsecase) DeleteOne(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err := use.repo.DeleteOne(ctx, int32(id))
	if err != nil {
		return err
	}

	return nil
}

func (use *FatUsecase) GetByID(id int) (model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByID(ctx, int32(id))
	if err != nil {
		return model.Fat{}, err
	}

	return (model.Fat)(res), nil
}

func (use *FatUsecase) FindByStates(state string, pag dto.Pagination) ([]model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByStates(ctx, state, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var fats []model.Fat
	for _, e := range res {
		fats = append(fats, (model.Fat)(e))
	}

	return fats, nil
}

func (use *FatUsecase) FindByMunicipality(state, municipality string, pag dto.Pagination) ([]model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByMunicipality(ctx, state, municipality, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var fats []model.Fat
	for _, e := range res {
		fats = append(fats, (model.Fat)(e))
	}

	return fats, nil
}

func (use *FatUsecase) FindByCounty(state, municipality, county string, pag dto.Pagination) ([]model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByCounty(ctx, state, municipality, county, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var fats []model.Fat
	for _, e := range res {
		fats = append(fats, (model.Fat)(e))
	}

	return fats, nil
}

func (use *FatUsecase) FindBytOdn(state, municipality, county, odn string, pag dto.Pagination) ([]model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindBytOdn(ctx, state, municipality, county, odn, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var fats []model.Fat
	for _, e := range res {
		fats = append(fats, (model.Fat)(e))
	}

	return fats, nil
}
