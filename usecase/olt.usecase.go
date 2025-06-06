package usecase

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/internal/snmp"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/metalpoch/ultra-monitor/repository"
)

type OltUsecase struct {
	repo repository.OltRepository
}

func NewOltUsecase(db *sqlx.DB) *OltUsecase {
	return &OltUsecase{repository.NewOltRepository(db)}
}

func (uc *OltUsecase) Add(olt dto.NewOlt) error {
	info, err := snmp.NewSnmp(snmp.Config{
		IP:        olt.IP,
		Community: olt.Community,
		Timeout:   5 * time.Second,
	}).OltSysQuery()
	if err != nil {
		return err
	}

	return uc.repo.Add(context.Background(), &entity.Olt{
		IP:          olt.IP,
		SysName:     info.SysName,
		SysLocation: info.SysLocation,
		Community:   olt.Community,
		IsAlive:     true,
		LastCheck:   time.Now(),
	})
}

func (uc *OltUsecase) Olt(ip string) (model.Olt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e, err := uc.repo.Olt(ctx, ip)
	return (model.Olt)(e), err
}

func (uc *OltUsecase) Olts() ([]model.Olt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.Olts(ctx)
	var olts []model.Olt
	for _, e := range res {
		olts = append(olts, (model.Olt)(e))
	}

	return olts, err
}

func (uc *OltUsecase) Delete(ip string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return uc.repo.Delete(ctx, ip)

}

func (uc *OltUsecase) GetAllIP() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return uc.repo.GetAllIP(ctx)
}

func (uc *OltUsecase) GetAllSysname() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return uc.repo.GetAllSysname(ctx)
}

func (uc *OltUsecase) OltsByState(state string) ([]model.Olt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.OltsByState(ctx, state)
	if err != nil {
		return nil, err
	}

	var olts []model.Olt
	for _, e := range res {
		olts = append(olts, (model.Olt)(e))
	}

	return olts, err
}

func (uc *OltUsecase) OltsByMunicipality(state, municipality string) ([]model.Olt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.OltsByMunicipality(ctx, state, municipality)
	if err != nil {
		return nil, err
	}

	var olts []model.Olt
	for _, e := range res {
		olts = append(olts, (model.Olt)(e))
	}

	return olts, err
}

func (uc *OltUsecase) OltsByCounty(state, municipality, county string) ([]model.Olt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := uc.repo.OltsByCounty(ctx, state, municipality, county)
	if err != nil {
		return nil, err
	}

	var olts []model.Olt
	for _, e := range res {
		olts = append(olts, (model.Olt)(e))
	}

	return olts, err
}
