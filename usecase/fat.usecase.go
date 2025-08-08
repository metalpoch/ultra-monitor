package usecase

import (
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/repository"
)

type FatUsecase struct {
	repo repository.FatRepository
}

func NewFatUsecase(db *sqlx.DB) *FatUsecase {
	return &FatUsecase{repository.NewFatRepository(db)}
}

func (use *FatUsecase) GetAll(pag dto.Pagination) ([]dto.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.AllInfo(ctx, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var taos []dto.Fat
	for _, t := range res {
		taos = append(taos, (dto.Fat)(t))
	}

	return taos, nil
}

func (use *FatUsecase) UpsertFats(file multipart.File, date time.Time) (int64, error) {
	tmpFile, err := os.CreateTemp("", "*.csv")
	if err != nil {
		return 0, err
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, file); err != nil {
		return 0, err
	}

	cmd := exec.Command("./fats-csv-to-json", tmpFile.Name())
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var data []dto.UpsertFat
	if err := json.Unmarshal(output, &data); err != nil {
		return 0, err
	}

	var fats []entity.UpsertFat
	for _, d := range data {
		fat := (entity.UpsertFat)(d)
		fat.Date = date
		fats = append(fats, fat)
	}
	return use.repo.UpsertFats(context.Background(), fats)
}

func (use *FatUsecase) GetByID(id int) (dto.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByID(ctx, int32(id))
	if err != nil {
		return dto.Fat{}, err
	}

	return (dto.Fat)(res), nil
}

func (use *FatUsecase) GetRegions() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetRegions(ctx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (use *FatUsecase) GetStates(region string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetStates(ctx, region)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (use *FatUsecase) GetMunicipalities(region, state string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetMunicipalities(ctx, region, state)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (use *FatUsecase) GetCounties(region, state, municipality string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetCounties(ctx, region, state, municipality)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (use *FatUsecase) GetODN(region, state, municipality, county string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetODN(ctx, region, state, municipality, county)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (use *FatUsecase) GetFat(region, state, municipality, county, odn string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetFat(ctx, region, state, municipality, county, odn)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (use *FatUsecase) FindByStates(state string, pag dto.Pagination) ([]dto.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByStates(ctx, state, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var fats []dto.Fat
	for _, e := range res {
		fats = append(fats, (dto.Fat)(e))
	}

	return fats, nil
}

func (use *FatUsecase) FindByMunicipality(state, municipality string, pag dto.Pagination) ([]dto.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByMunicipality(ctx, state, municipality, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var fats []dto.Fat
	for _, e := range res {
		fats = append(fats, (dto.Fat)(e))
	}

	return fats, nil
}

func (use *FatUsecase) FindByCounty(state, municipality, county string, pag dto.Pagination) ([]dto.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByCounty(ctx, state, municipality, county, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var fats []dto.Fat
	for _, e := range res {
		fats = append(fats, (dto.Fat)(e))
	}

	return fats, nil
}

func (use *FatUsecase) FindBytOdn(state, municipality, county, odn string, pag dto.Pagination) ([]dto.Fat, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.FindByOdn(ctx, state, municipality, county, odn, pag.Page, pag.Limit)
	if err != nil {
		return nil, err
	}

	var fats []dto.Fat
	for _, e := range res {
		fats = append(fats, (dto.Fat)(e))
	}

	return fats, nil
}

func (use *FatUsecase) GetAllFatStatus() ([]dto.FatStatusSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetAllFatStatus(ctx)
	if err != nil {
		return nil, err
	}

	var fats []dto.FatStatusSummary
	for _, e := range res {
		fats = append(fats, (dto.FatStatusSummary)(e))
	}

	return fats, nil
}

func (use *FatUsecase) GetAllFatStatusByRegion(region string) ([]dto.FatStatusSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetAllFatStatusByRegion(ctx, region)
	if err != nil {
		return nil, err
	}

	var fats []dto.FatStatusSummary
	for _, e := range res {
		fats = append(fats, (dto.FatStatusSummary)(e))
	}

	return fats, nil
}

func (use *FatUsecase) GetAllFatStatusByState(state string) ([]dto.FatStatusSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetAllFatStatusByState(ctx, state)
	if err != nil {
		return nil, err
	}

	var fats []dto.FatStatusSummary
	for _, e := range res {
		fats = append(fats, (dto.FatStatusSummary)(e))
	}

	return fats, nil
}

func (use *FatUsecase) GetAllFatStatusByMunicipality(state, municipality string) ([]dto.FatStatusSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetAllFatStatusByMunicipality(ctx, state, municipality)
	if err != nil {
		return nil, err
	}

	var fats []dto.FatStatusSummary
	for _, e := range res {
		fats = append(fats, (dto.FatStatusSummary)(e))
	}

	return fats, nil
}

func (use *FatUsecase) GetAllFatStatusByCounty(state, municipality, county string) ([]dto.FatStatusSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetAllFatStatusByCounty(ctx, state, municipality, county)
	if err != nil {
		return nil, err
	}

	var fats []dto.FatStatusSummary
	for _, e := range res {
		fats = append(fats, (dto.FatStatusSummary)(e))
	}

	return fats, nil
}

func (use *FatUsecase) GetAllFatStatusByODN(state, municipality, county, odn string) ([]dto.FatStatusSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetAllFatStatusByODN(ctx, state, municipality, county, odn)
	if err != nil {
		return nil, err
	}

	var fats []dto.FatStatusSummary
	for _, e := range res {
		fats = append(fats, (dto.FatStatusSummary)(e))
	}

	return fats, nil
}

func (use *FatUsecase) GetAllFatStatusByFat(state, municipality, county, odn, fat string) ([]dto.FatStatusSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetAllFatStatusByFat(ctx, state, municipality, county, odn, fat)
	if err != nil {
		return nil, err
	}

	var fats []dto.FatStatusSummary
	for _, e := range res {
		fats = append(fats, (dto.FatStatusSummary)(e))
	}

	return fats, nil
}
