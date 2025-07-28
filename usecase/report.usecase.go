package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/dto"
	"github.com/metalpoch/ultra-monitor/model"
	"github.com/metalpoch/ultra-monitor/repository"
)

type ReportUsecase struct {
	repo repository.ReportRepository
}

func NewReportUsecase(db *sqlx.DB) *ReportUsecase {
	return &ReportUsecase{repository.NewReportRepository(db)}
}

func (use *ReportUsecase) Add(rp dto.NewReport) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	newReport := &entity.Report{
		UserID:           rp.UserID,
		Category:         rp.Category,
		Basepath:         rp.Basepath,
		ContentType:      rp.File.Header.Get("Content-Type"),
		OriginalFilename: rp.File.Filename,
	}
	err := use.repo.Add(ctx, newReport)
	if err != nil {
		return "", err
	}
	return newReport.ID.String(), err
}

func (use *ReportUsecase) Get(id string) (*model.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fileID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	res, err := use.repo.Get(ctx, fileID)
	if err != nil {
		fmt.Println("ERROR", err.Error())
		return nil, err
	}
	return (*model.Report)(res), nil
}

func (use *ReportUsecase) GetReportsByUser(userID int, page dto.Pagination) ([]model.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetReportsByUser(ctx, int32(userID), page.Page, page.Limit)
	if err != nil {
		return nil, err
	}

	var reports []model.Report
	for _, rp := range res {
		reports = append(reports, (model.Report)(rp))
	}
	return reports, err
}

func (use *ReportUsecase) GetReportsByCategory(category string, page dto.Pagination) ([]model.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetReportsByCategory(ctx, category, page.Page, page.Limit)
	if err != nil {
		return nil, err
	}

	var reports []model.Report
	for _, rp := range res {
		reports = append(reports, (model.Report)(rp))
	}
	return reports, err
}

func (use *ReportUsecase) GetCategories() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	res, err := use.repo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (use *ReportUsecase) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fileID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return use.repo.Delete(ctx, fileID)
}
