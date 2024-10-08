package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/metalpoch/olt-blueprint/common/constants"
	"github.com/metalpoch/olt-blueprint/common/entity"
	"github.com/metalpoch/olt-blueprint/common/model"
	"github.com/metalpoch/olt-blueprint/common/pkg/tracking"
	"github.com/metalpoch/olt-blueprint/report/repository"
	"gorm.io/gorm"
)

type reportUsecase struct {
	repo     repository.ReportRepository
	telegram tracking.Telegram
}

func NewReportUsecase(db *gorm.DB, telegram tracking.Telegram) *reportUsecase {
	return &reportUsecase{repository.NewReportRepository(db), telegram}
}

func (use reportUsecase) Add(rp *model.NewReport) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	newReport := &entity.Report{
		UserID:   rp.UserID,
		Category: rp.Category,
		Filename: rp.Filename,
		Filepath: constants.BASE_FILEPATH,
	}
	err := use.repo.Add(ctx, newReport)
	if err != nil {
		go use.telegram.Notification(
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
		go use.telegram.Notification(
			constants.MODULE_REPORT,
			constants.CATEGORY_DATABASE,
			fmt.Sprintf("(feedUsecase).GetDevice - use.repo.GetDevice(ctx, %s)", id),
			err,
		)
		return nil, err
	}
	return (*model.Report)(res), nil
}
