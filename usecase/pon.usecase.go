package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/metalpoch/ultra-monitor/repository"
)

type PonUsecase struct {
	repo repository.PonRepository
}

func NewPonUsecase(db *sqlx.DB) *PonUsecase {
	return &PonUsecase{repository.NewPonRepository(db)}
}

func (uc *PonUsecase) GetAllByDevice(sysname string) ([]model.Pon, error) {
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

func (uc *PonUsecase) PonByOltAndPort(sysname string, shell, card, port int) (*model.Pon, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	ifname := fmt.Sprintf("GPON %d/%d/%d", shell, card, port)
	pon, err := uc.repo.PonByPort(ctx, sysname, ifname)
	if err != nil {
		return nil, err
	}
	return (*model.Pon)(&pon), nil
}
