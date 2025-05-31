package repository

import (
	"context"
	"path"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
	"github.com/metalpoch/ultra-monitor/internal/constants"
)

type ReportRepository interface {
	Add(ctx context.Context, report *entity.Report) error
	Get(ctx context.Context, id uuid.UUID) (*entity.Report, error)
	GetCategories(ctx context.Context) ([]string, error)
	GetReportsByUser(ctx context.Context, userID int32, page, limit uint16) ([]entity.Report, error)
	GetReportsByCategory(ctx context.Context, category string, page, limit uint16) ([]entity.Report, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type reportRepository struct {
	db *sqlx.DB
}

func NewReportRepository(db *sqlx.DB) *reportRepository {
	return &reportRepository{db}
}

func (repo *reportRepository) Add(ctx context.Context, report *entity.Report) error {
	report.ID = uuid.New()
	report.Filepath = path.Join(report.Basepath, report.ID.String())
	_, err := repo.db.NamedExecContext(ctx, constants.SQL_INSERT_REPORT, report)
	return err
}

func (repo *reportRepository) Get(ctx context.Context, id uuid.UUID) (*entity.Report, error) {
	var report entity.Report
	err := repo.db.GetContext(ctx, &report, constants.SQL_SELECT_REPORT_BY_ID, id)
	return &report, err
}

func (repo *reportRepository) GetCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := repo.db.SelectContext(ctx, &categories, constants.SQL_SELECT_DISTINCT_REPORT_CATEGORIES)
	return categories, err
}

func (repo *reportRepository) GetReportsByUser(ctx context.Context, userID int32, page, limit uint16) ([]entity.Report, error) {
	offset := (page - 1) * limit
	var reports []entity.Report
	err := repo.db.SelectContext(ctx, &reports, constants.SQL_SELECT_REPORTS_BY_USER, userID, limit, offset)
	return reports, err
}

func (repo *reportRepository) GetReportsByCategory(ctx context.Context, category string, page, limit uint16) ([]entity.Report, error) {
	offset := (page - 1) * limit
	var reports []entity.Report
	err := repo.db.SelectContext(ctx, &reports, constants.SQL_SELECT_REPORTS_BY_CATEGORY, category, limit, offset)
	return reports, err
}

func (repo *reportRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := repo.db.ExecContext(ctx, constants.SQL_DELETE_REPORT_BY_ID, id)
	return err
}
