package repository

import (
	"context"
	"path"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/metalpoch/ultra-monitor/entity"
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
	query := `
    INSERT INTO reports (
        id, category, original_filename, content_type, basepath, filepath, user_id
    ) VALUES (
        :id, :category, :original_filename, :content_type, :basepath, :filepath, :user_id
    )`
	report.ID = uuid.New()
	report.Filepath = path.Join(report.Basepath, report.ID.String())
	_, err := repo.db.NamedExecContext(ctx, query, report)
	return err
}

func (repo *reportRepository) Get(ctx context.Context, id uuid.UUID) (*entity.Report, error) {
	var report entity.Report
	query := `SELECT * FROM reports WHERE id = $1 ORDER BY created_at`
	err := repo.db.GetContext(ctx, &report, query, id)
	return &report, err
}

func (repo *reportRepository) GetCategories(ctx context.Context) ([]string, error) {
	var categories []string
	query := `SELECT DISTINCT category FROM reports ORDER BY category`
	err := repo.db.SelectContext(ctx, &categories, query)
	return categories, err
}

func (repo *reportRepository) GetReportsByUser(ctx context.Context, userID int32, page, limit uint16) ([]entity.Report, error) {
	offset := (page - 1) * limit
	var reports []entity.Report
	query := `SELECT * FROM reports WHERE user_id = $1 LIMIT $2 OFFSET $3 ORDER BY created_at`
	err := repo.db.SelectContext(ctx, &reports, query, userID, limit, offset)
	return reports, err
}

func (repo *reportRepository) GetReportsByCategory(ctx context.Context, category string, page, limit uint16) ([]entity.Report, error) {
	offset := (page - 1) * limit
	var reports []entity.Report
	query := `SELECT * FROM reports WHERE category = $1 LIMIT $2 OFFSET $3 ORDER BY created_at`
	err := repo.db.SelectContext(ctx, &reports, query, category, limit, offset)
	return reports, err
}

func (repo *reportRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM reports WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}
