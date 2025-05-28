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

type OltUsecase struct {
	repo repository.OltRepository
}

func NewOltUsecase(db *sqlx.DB) *OltUsecase {
	return &OltUsecase{repository.NewOltRepository(db)}
}

func (uc OltUsecase) Add(olt model.Olt) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	newDevice := &entity.Olt{
		IP:          olt.IP,
		SysName:     olt.SysName,
		SysLocation: olt.SysLocation,
		Community:   olt.Community,
		IsAlive:     olt.IsAlive,
		LastCheck:   time.Now(),
	}

	return uc.repo.Add(ctx, newDevice)
}

func (uc OltUsecase) Update(id uint64, olt dto.NewOlt) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e, err := uc.repo.Olt(ctx, id)
	if err != nil {
		return err
	}

	if olt.IP != "" {
		e.IP = olt.IP
	}
	if olt.Community != "" {
		e.Community = olt.Community
	}
	return uc.repo.Update(ctx, e)
}

func (uc OltUsecase) Delete(id uint64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return uc.repo.Delete(ctx, id)

}

func (uc OltUsecase) Olt(id uint64) (model.Olt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e, err := uc.repo.Olt(ctx, id)
	return (model.Olt)(e), err
}

func (uc OltUsecase) Olts(page, limit uint16) ([]model.Olt, error) {
	var devices []model.Olt
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.Olts(ctx, page, limit)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		devices = append(devices, (model.Olt)(e))
	}

	return devices, err
}

func (uc OltUsecase) OltsByState(state string, page, limit uint16) ([]model.Olt, error) {
	var devices []model.Olt
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.OltsByState(ctx, state, page, limit)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		devices = append(devices, (model.Olt)(e))
	}

	return devices, err
}

func (uc OltUsecase) OltsByCounty(state, county string, page, limit uint16) ([]model.Olt, error) {
	var devices []model.Olt
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.OltsByCounty(ctx, state, county, page, limit)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		devices = append(devices, (model.Olt)(e))
	}

	return devices, err
}
func (uc OltUsecase) OltsByMunicipality(state, county, municipality string, page, limit uint16) ([]model.Olt, error) {
	var devices []model.Olt
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.OltsByMunicipality(ctx, state, county, municipality, page, limit)
	if err != nil {
		return nil, err
	}

	for _, e := range res {
		devices = append(devices, (model.Olt)(e))
	}

	return devices, err
}
