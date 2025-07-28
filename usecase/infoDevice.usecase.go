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

type InfoDeviceUsecase struct {
	repo repository.InfoDeviceRepository
}

func NewInfoDeviceUsecase(db *sqlx.DB) *InfoDeviceUsecase {
	return &InfoDeviceUsecase{repository.NewInfoDeviceRepository(db)}
}

func (use *InfoDeviceUsecase) GetAll(pag dto.Pagination) ([]model.InfoDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.AllInfo(ctx, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var taos []model.InfoDevice
	for _, t := range res {
		taos = append(taos, (model.InfoDevice)(t))
	}

	return taos, nil
}

func (use *InfoDeviceUsecase) AddInfo(info dto.InfoDevice) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	_, err := use.repo.AddInfo(ctx, entity.InfoDevice{
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

func (use *InfoDeviceUsecase) DeleteOne(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	err := use.repo.DeleteOne(ctx, int32(id))
	if err != nil {
		return err
	}

	return nil
}

func (use *InfoDeviceUsecase) GetByID(id int) (model.InfoDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByID(ctx, int32(id))
	if err != nil {
		return model.InfoDevice{}, err
	}

	return (model.InfoDevice)(res), nil
}

func (use *InfoDeviceUsecase) FindByStates(state string, pag dto.Pagination) ([]model.InfoDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByStates(ctx, state, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var fats []model.InfoDevice
	for _, e := range res {
		fats = append(fats, (model.InfoDevice)(e))
	}

	return fats, nil
}

func (use *InfoDeviceUsecase) FindByMunicipality(state, municipality string, pag dto.Pagination) ([]model.InfoDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByMunicipality(ctx, state, municipality, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var fats []model.InfoDevice
	for _, e := range res {
		fats = append(fats, (model.InfoDevice)(e))
	}

	return fats, nil
}

func (use *InfoDeviceUsecase) FindByCounty(state, municipality, county string, pag dto.Pagination) ([]model.InfoDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByCounty(ctx, state, municipality, county, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var fats []model.InfoDevice
	for _, e := range res {
		fats = append(fats, (model.InfoDevice)(e))
	}

	return fats, nil
}

func (use *InfoDeviceUsecase) FindBytOdn(state, municipality, county, odn string, pag dto.Pagination) ([]model.InfoDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindBytOdn(ctx, state, municipality, county, odn, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var fats []model.InfoDevice
	for _, e := range res {
		fats = append(fats, (model.InfoDevice)(e))
	}

	return fats, nil
}
