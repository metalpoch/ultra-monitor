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

func (use *FatUsecase) AddFat(fat model.Fat) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	entityFat := (entity.Fat)(fat)
	err := use.repo.AddFat(ctx, &entityFat)
	if err != nil {
		return err
	}

	return nil
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

	res, err := use.repo.GetByID(ctx, int32(id))
	if err != nil {
		return model.Fat{}, err
	}

	return (model.Fat)(res), nil
}

func (use *FatUsecase) GetTraffic(id int, dates dto.RangeDate) ([]model.Traffic, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetTraffic(ctx, id, dates.InitDate, dates.EndDate)
	if err != nil {
		return nil, err
	}

	var traffic []model.Traffic
	for _, t := range res {
		traffic = append(traffic, (model.Traffic)(t))
	}

	return traffic, err
}

func (use *FatUsecase) GetStates() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetStates(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (use *FatUsecase) GetMunicipality(state string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetMunicipality(ctx, state)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (use *FatUsecase) GetCounty(state, municipality string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetCounty(ctx, state, municipality)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (use *FatUsecase) GetOdn(state, municipality, county string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetOdn(ctx, state, municipality, county)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (use *FatUsecase) GetFatsByStates(state string) ([]model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetFatsByStates(ctx, state)
	if err != nil {
		return nil, err
	}

	var fats []model.Fat
	for _, e := range res {
		fats = append(fats, (model.Fat)(e))
	}

	return fats, nil
}

func (use *FatUsecase) GetFatsByMunicipality(state, municipality string) ([]model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetFatsByMunicipality(ctx, state, municipality)
	if err != nil {
		return nil, err
	}

	var fats []model.Fat
	for _, e := range res {
		fats = append(fats, (model.Fat)(e))
	}

	return fats, nil
}

func (use *FatUsecase) GetFatsByCounty(state, municipality, county string) ([]model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetFatsByCounty(ctx, state, municipality, county)
	if err != nil {
		return nil, err
	}

	var fats []model.Fat
	for _, e := range res {
		fats = append(fats, (model.Fat)(e))
	}

	return fats, nil
}

func (use *FatUsecase) GetFatsBytOdn(state, municipality, county, odn string) ([]model.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetFatsBytOdn(ctx, state, municipality, county, odn)
	if err != nil {
		return nil, err
	}

	var fats []model.Fat
	for _, e := range res {
		fats = append(fats, (model.Fat)(e))
	}

	return fats, nil
}
