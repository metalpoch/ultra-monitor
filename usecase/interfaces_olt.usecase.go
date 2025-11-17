package usecase

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/repository"
)

type InterfaceOltUsecase struct {
	repo repository.InterfaceOltRepository
}

func NewInterfaceOltUsecase(db *sqlx.DB) *InterfaceOltUsecase {
	return &InterfaceOltUsecase{
		repo: repository.NewInterfaceOltRepository(db),
	}
}

func (u *InterfaceOltUsecase) GetAll() ([]dto.InterfacesDetailedOlt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := u.repo.GetAll(ctx)

	var olts []dto.InterfacesDetailedOlt

	for _, e := range res {
		olts = append(olts, (dto.InterfacesDetailedOlt)(e))
	}

	return olts, err
}

func (u *InterfaceOltUsecase) Update(data dto.InterfacesOlt) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.repo.Update(ctx, data.OltIP, data.OltVerbose)
}
