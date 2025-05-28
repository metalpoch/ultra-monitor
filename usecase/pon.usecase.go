package usecase

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/olt-blueprint/model"
	"github.com/metalpoch/olt-blueprint/repository"
)

type PonUsecase struct {
	repo repository.PonRepository
}

func NewPonUsecase(db *sqlx.DB) *PonUsecase {
	return &PonUsecase{repository.NewPonRepository(db)}
}

func (uc PonUsecase) GetAllByDevice(sysname string) ([]model.Pon, error) {
	var interfaces []model.Pon
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := uc.repo.PonsByOLT(ctx, sysname)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		interfaces = append(interfaces, (model.Pon)(e))
	}

	return interfaces, nil
}

func (uc PonUsecase) PonByOltAndPort(sysname, port string) (*model.Pon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pon, err := uc.repo.PonByPort(ctx, sysname, port)
	if err != nil {
		return nil, err
	}
	return (*model.Pon)(&pon), nil
}
