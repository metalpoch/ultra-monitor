package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/report/repository"
	"gorm.io/gorm"
)

type reportUsecase struct {
	repo     repository.ReportRepository
	telegram tracking.SmartModule
}

func NewReportUsecase(db *gorm.DB, telegram tracking.SmartModule) *reportUsecase {
	return &reportUsecase{repository.NewReportRepository(db), telegram}
}

func (use reportUsecase) Add(rp *model.NewReport) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newReport := &entity.Report{
		UserID:           rp.UserID,
		Category:         rp.Category,
		ContentType:      rp.File.Header.Get("Content-Type"),
		OriginalFilename: rp.File.Filename,
	}
	err := use.repo.Add(ctx, newReport)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(reportUsecase).Add - use.repo.Add(ctx, %v)", newReport),
			err,
		)
	}

	return newReport.ID.String(), err
}

func (use reportUsecase) Get(id string) (*model.Report, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.Get(ctx, id)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(reportUsecase).Get - use.repo.Get(ctx, %s)", id),
			err,
		)
		return nil, err
	}
	return (*model.Report)(res), nil
}

func (use reportUsecase) GetReports(query *model.FindReports) ([]*model.ReportResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetReports(ctx, query.Category, query.UserID)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			"(reportUsecase).GetReports - use.repo.GetReports(ctx, %s)",
			err,
		)
	}

	var reports []*model.ReportResponse
	for _, e := range res {
		reports = append(reports, &model.ReportResponse{
			Category:         e.Category,
			OriginalFilename: e.OriginalFilename,
			ContentType:      e.ContentType,
			Filepath:         e.Filepath,
			CreatedAt:        e.CreatedAt,
			UpdatedAt:        e.UpdatedAt,
			DeletedAt:        e.DeletedAt.Time,
			User: model.UserLite{
				ID:       e.User.ID,
				Email:    e.User.Email,
				Fullname: e.User.Fullname,
			},
		})
	}
	return reports, err
}

func (use reportUsecase) GetCategories() ([]*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := use.repo.GetCategories(ctx)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			"(reportUsecase).GetCategories - use.repo.GetCategories(ctx, %s)",
			err,
		)
	}

	return res, err
}

func (use reportUsecase) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(reportUsecase).Get - use.repo.Get(ctx, %s)", id),
			err,
		)
	}

	err = use.repo.Delete(ctx, parsedUUID)
	if err != nil {
		go use.telegram.SendMessage(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(reportUsecase).Get - use.repo.Delete(ctx, %s)", parsedUUID),
			err,
		)
	}

	return err
}
